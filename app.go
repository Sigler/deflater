package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"deflater/internal/appx"
	"deflater/internal/catalog"
	"deflater/internal/config"
	"deflater/internal/elevate"
	"deflater/internal/engine"
	"deflater/internal/logging"
	"deflater/internal/reg"
	"deflater/internal/schtask"
	"deflater/internal/toast"
	"deflater/internal/update"
)

// App is the bridge the frontend calls. Every exported method becomes a
// typed TypeScript function in frontend/wailsjs.
type App struct {
	ctx context.Context
	eng *engine.Engine
	// mu serializes config read-modify-write cycles within this process;
	// config.Update adds the cross-process lock on top.
	mu sync.Mutex
	// applying is true while Apply is mutating the machine, so the window
	// cannot be closed mid-apply.
	applying atomic.Bool
	// dirty mirrors the frontend's count of staged-but-unapplied changes
	// so closing the window can warn before losing them.
	dirty atomic.Int32
}

func NewApp() *App {
	return &App{eng: &engine.Engine{Appx: &appx.Service{}}}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	_ = toast.Register("")
	logging.Logf("app: started (elevated=%v)", elevate.IsElevated())
}

// RegOpInfo mirrors reg.Op for the frontend, so the generated bindings
// stay self-contained in this package.
type RegOpInfo struct {
	Hive   string `json:"hive"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	Value  uint32 `json:"value"`
	Revert string `json:"revert"`
}

// FixState is one catalog entry plus its live status on this machine.
type FixState struct {
	ID       string      `json:"id"`
	Category string      `json:"category"`
	Kind     string      `json:"kind"`
	Caution  bool        `json:"caution"`
	Profiles []string    `json:"profiles"`
	Reg      []RegOpInfo `json:"reg,omitempty"`
	Appx     []string    `json:"appx,omitempty"`
	Status   string      `json:"status"`
}

// Report is everything the UI needs to render.
type Report struct {
	Version     string         `json:"version"`
	Elevated    bool           `json:"elevated"`
	Categories  []string       `json:"categories"`
	Fixes       []FixState     `json:"fixes"`
	Managed     []string       `json:"managed"`
	Maintenance bool           `json:"maintenance"`
	Watcher     bool           `json:"watcher"`
	Alerts      []config.Alert `json:"alerts"`
	// TaskMismatch is true when maintenance/watcher is on in config but
	// the scheduled task is missing (or vice versa) — a self-heal hint.
	TaskMismatch bool `json:"taskMismatch"`
	// ConflictingTasks are scheduled tasks from earlier debloat tools that
	// fight Deflater's state; the UI offers to remove them.
	ConflictingTasks []schtask.ForeignTask `json:"conflictingTasks"`
	// Pending is set when an apply was interrupted by an elevation
	// relaunch; the UI resumes it after confirmation.
	Pending *config.Pending `json:"pending"`
}

// GetReport gathers catalog, statuses, and settings. It re-reads the
// installed-app list so each report reflects current reality (and a
// transient enumeration failure recovers on the next call).
func (a *App) GetReport() (Report, error) {
	_ = a.eng.Appx.Refresh()
	cfg := config.Load()
	r := Report{
		Version:     appVersion,
		Elevated:    elevate.IsElevated(),
		Categories:  catalog.Categories,
		Managed:     cfg.Managed,
		Maintenance: cfg.Maintenance,
		Watcher:     cfg.WatcherEnabled,
		Alerts:      cfg.Alerts,
		Pending:     cfg.Pending,
	}
	if r.Managed == nil {
		r.Managed = []string{}
	}
	if r.Alerts == nil {
		r.Alerts = []config.Alert{}
	}
	wantTask := cfg.Maintenance || cfg.WatcherEnabled
	r.TaskMismatch = wantTask != schtask.IsInstalled()
	r.ConflictingTasks = schtask.DetectForeign()
	for _, f := range catalog.Fixes() {
		status, err := a.eng.Status(f)
		if err != nil {
			logging.Logf("report: %s status failed: %v", f.ID, err)
			status = "unknown"
		}
		state := FixState{
			ID:       f.ID,
			Category: f.Category,
			Kind:     string(f.Kind),
			Caution:  f.Caution,
			Profiles: f.Profiles,
			Appx:     f.Appx,
			Status:   status,
		}
		for _, op := range f.Reg {
			state.Reg = append(state.Reg, RegOpInfo(op))
		}
		r.Fixes = append(r.Fixes, state)
	}
	return r, nil
}

// FixResult reports one fix's outcome from an apply run.
type FixResult struct {
	ID     string `json:"id"`
	OK     bool   `json:"ok"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
	// Phase is "start" for the pre-work event and "done" for the result.
	Phase string `json:"phase"`
}

// ApplyOutcome is the result of an Apply call. If NeedsElevation is set,
// nothing was changed; the UI should confirm and call SaveAndElevate.
type ApplyOutcome struct {
	NeedsElevation bool        `json:"needsElevation"`
	Results        []FixResult `json:"results"`
	// SaveWarning is set when the machine changes applied but persisting
	// the record failed; the changes are real but may not be tracked.
	SaveWarning string `json:"saveWarning,omitempty"`
}

// Apply enables and disables fixes. Enable ids are applied (and their
// pre-apply registry state captured for exact revert), disable ids
// reverted. Progress events stream as "apply:progress". Unelevated it
// changes nothing and asks for an elevated relaunch.
func (a *App) Apply(enable []string, disable []string) (ApplyOutcome, error) {
	if !elevate.IsElevated() {
		return ApplyOutcome{NeedsElevation: true}, nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.applying.Store(true)
	defer a.applying.Store(false)

	_ = a.eng.Appx.Refresh()
	existing := config.Load()

	// A change records how the persisted config should reflect one fix's
	// outcome; only successful steps produce one.
	type change struct {
		id       string
		managed  bool
		setSnaps bool
		snaps    []reg.Snapshot
	}
	var results []FixResult
	var changes []change

	emit := func(res FixResult) {
		if a.ctx != nil {
			wruntime.EventsEmit(a.ctx, "apply:progress", res)
		}
	}

	for _, id := range enable {
		fix, ok := catalog.ByID(id)
		if !ok {
			continue
		}
		emit(FixResult{ID: id, Phase: "start"})
		res := FixResult{ID: id, OK: true, Phase: "done"}
		// Capture the prior state once, before the first apply, so revert
		// can restore what the user actually had.
		var snaps []reg.Snapshot
		capture := false
		if _, have := existing.Snapshots[id]; !have {
			if s, err := a.eng.Capture(fix); err == nil {
				snaps, capture = s, true
			}
		}
		if err := a.eng.Apply(fix, true); err != nil {
			res.OK = false
			res.Error = err.Error()
			logging.Logf("apply: %s failed: %v", id, err)
		} else {
			changes = append(changes, change{id: id, managed: true, setSnaps: capture, snaps: snaps})
		}
		res.Status, _ = a.eng.Status(fix)
		results = append(results, res)
		emit(res)
	}
	for _, id := range disable {
		fix, ok := catalog.ByID(id)
		if !ok {
			continue
		}
		emit(FixResult{ID: id, Phase: "start"})
		res := FixResult{ID: id, OK: true, Phase: "done"}
		if err := a.eng.Revert(fix, existing.Snapshots[id]); err != nil {
			res.OK = false
			res.Error = err.Error()
			logging.Logf("revert: %s failed: %v", id, err)
		} else {
			changes = append(changes, change{id: id, managed: false, setSnaps: true, snaps: nil})
		}
		res.Status, _ = a.eng.Status(fix)
		results = append(results, res)
		emit(res)
	}

	// Persist under the cross-process lock, merging only the fields this
	// apply owns, then sync the scheduled task while we hold admin rights.
	saveErr := config.Update(func(c *config.Config) error {
		for _, ch := range changes {
			c.SetManaged(ch.id, ch.managed)
			if ch.setSnaps {
				c.SetSnapshot(ch.id, ch.snaps)
			}
		}
		c.Pending = nil
		return a.syncTask(c)
	})

	out := ApplyOutcome{Results: results}
	if saveErr != nil {
		// The machine changed; don't discard the results just because
		// recording them failed. Surface it as a warning instead.
		logging.Logf("apply: persist failed: %v", saveErr)
		out.SaveWarning = saveErr.Error()
	}
	return out, nil
}

// SaveAndElevate stores the confirmed apply request (tagged so the
// elevated instance can verify it), relaunches with admin rights, and
// closes this instance.
func (a *App) SaveAndElevate(enable []string, disable []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	token := newToken()
	err := config.Update(func(c *config.Config) error {
		c.Pending = &config.Pending{
			Enable:  enable,
			Disable: disable,
			Token:   token,
			Created: time.Now().Format(time.RFC3339),
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := elevate.Relaunch(); err != nil {
		// UAC declined; forget the request so it cannot fire later.
		_ = config.Update(func(c *config.Config) error { c.Pending = nil; return nil })
		logging.Logf("elevate: relaunch failed or was declined: %v", err)
		return fmt.Errorf("Windows did not grant administrator rights")
	}
	// The staged changes are being handed off, not abandoned: clear dirty
	// so the close does not warn about losing them.
	logging.Logf("elevate: saved pending (%d enable, %d disable), relaunching as admin",
		len(enable), len(disable))
	a.dirty.Store(0)
	wruntime.Quit(a.ctx)
	return nil
}

// TakePending atomically returns and clears a valid pending request for
// the elevated instance to resume, or nil if there is none, it is
// expired, or this instance is not elevated. Consume-on-read means a
// request can never fire twice.
func (a *App) TakePending() (*config.Pending, error) {
	if !elevate.IsElevated() {
		return nil, nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	var p *config.Pending
	err := config.Update(func(c *config.Config) error {
		if c.Pending == nil {
			return nil
		}
		if c.Pending.Expired(time.Now()) {
			logging.Logf("pending: discarding expired request")
			c.Pending = nil
			return nil
		}
		p = c.Pending
		c.Pending = nil
		return nil
	})
	if p != nil {
		logging.Logf("pending: resuming saved apply (%d enable, %d disable)", len(p.Enable), len(p.Disable))
	}
	return p, err
}

func newToken() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// RemoveConflictingTasks deletes foreign debloat-tool scheduled tasks by
// name. Only vetted names are honored (RemoveForeign enforces this).
// Must be elevated: these tasks run with admin rights. Called directly
// when already elevated, or on the elevated resume for a staged removal.
func (a *App) RemoveConflictingTasks(names []string) error {
	if !elevate.IsElevated() {
		return fmt.Errorf("administrator rights are required to remove a scheduled task")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	var firstErr error
	for _, n := range names {
		if err := schtask.RemoveForeign(n); err != nil {
			logging.Logf("foreign task: remove %q failed: %v", n, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		logging.Logf("foreign task: removed %q", n)
	}
	return firstErr
}

// StageTaskRemovalAndElevate saves a foreign-task removal request and
// relaunches with admin rights to carry it out, mirroring SaveAndElevate.
// The elevated instance picks it up via TakePending. Unknown task names
// are rejected before any relaunch.
func (a *App) StageTaskRemovalAndElevate(name string) error {
	if !schtask.IsKnownForeign(name) {
		return fmt.Errorf("unrecognized task %q", name)
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	token := newToken()
	err := config.Update(func(c *config.Config) error {
		c.Pending = &config.Pending{
			RemoveTasks: []string{name},
			Token:       token,
			Created:     time.Now().Format(time.RFC3339),
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := elevate.Relaunch(); err != nil {
		_ = config.Update(func(c *config.Config) error { c.Pending = nil; return nil })
		logging.Logf("elevate: relaunch for task removal declined: %v", err)
		return fmt.Errorf("Windows did not grant administrator rights")
	}
	logging.Logf("elevate: saved pending task removal %q, relaunching as admin", name)
	a.dirty.Store(0)
	wruntime.Quit(a.ctx)
	return nil
}

// SetMaintenance turns automatic re-applying on or off. Unelevated the
// preference saves and takes effect on the next elevated apply
// (NeedsElevation); elevated, a task-sync failure is returned as an
// error so the UI can roll its toggle back.
func (a *App) SetMaintenance(on bool) (ToggleResult, error) {
	return a.setFlag(on, func(c *config.Config, v bool) { c.Maintenance = v })
}

// SetWatcher toggles silent-install alerts. Independent of maintenance:
// either switch keeps the shared scheduled task alive.
func (a *App) SetWatcher(on bool) (ToggleResult, error) {
	return a.setFlag(on, func(c *config.Config, v bool) { c.WatcherEnabled = v })
}

// ToggleResult tells the UI what actually happened to a settings toggle.
type ToggleResult struct {
	// Saved is true when the preference was persisted.
	Saved bool `json:"saved"`
	// NeedsElevation is true when the task will be registered on the
	// next elevated apply (unelevated toggle).
	NeedsElevation bool `json:"needsElevation"`
}

func (a *App) setFlag(on bool, set func(*config.Config, bool)) (ToggleResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	elevated := elevate.IsElevated()
	var syncErr error
	err := config.Update(func(c *config.Config) error {
		set(c, on)
		if elevated {
			if e := a.syncTask(c); e != nil {
				set(c, !on) // roll the preference back to match reality
				syncErr = e
			}
		}
		return nil
	})
	if err != nil {
		return ToggleResult{}, err
	}
	if syncErr != nil {
		return ToggleResult{}, fmt.Errorf("couldn't update the scheduled task: %w", syncErr)
	}
	return ToggleResult{Saved: true, NeedsElevation: on && !elevated}, nil
}

// syncTask makes the real scheduled task match the config: it exists
// while either maintenance or the watcher wants it, and is removed when
// both are off. Must be called elevated. On install it first copies the
// exe to an admin-only per-machine location so the auto-elevating task
// cannot be pointed at a user-writable binary.
func (a *App) syncTask(cfg *config.Config) error {
	if cfg.Maintenance || cfg.WatcherEnabled {
		exePath, err := installSelf()
		if err == nil {
			err = schtask.Install(exePath)
		}
		if err != nil {
			logging.Logf("task: enable failed: %v", err)
			return err
		}
		logging.Logf("task: registered at %s", exePath)
		return nil
	}
	if err := schtask.Uninstall(); err != nil {
		logging.Logf("task: remove failed: %v", err)
		return err
	}
	return nil
}

// SetDirty records how many staged changes are awaiting apply.
func (a *App) SetDirty(n int) {
	a.dirty.Store(int32(n))
}

// beforeClose blocks closing mid-apply, and warns before discarding
// staged changes. Returning true prevents the close.
func (a *App) beforeClose(ctx context.Context) bool {
	if a.applying.Load() {
		wruntime.MessageDialog(ctx, wruntime.MessageDialogOptions{
			Type:    wruntime.InfoDialog,
			Title:   "Deflater",
			Message: "Deflater is still applying your changes. Please wait for it to finish.",
			Buttons: []string{"OK"},
		})
		return true
	}
	n := int(a.dirty.Load())
	if n <= 0 {
		return false
	}
	msg := fmt.Sprintf("You have %d pending changes that were never applied. Close without applying them?", n)
	if n == 1 {
		msg = "You have 1 pending change that was never applied. Close without applying it?"
	}
	choice, err := wruntime.MessageDialog(ctx, wruntime.MessageDialogOptions{
		Type:          wruntime.QuestionDialog,
		Title:         "Deflater",
		Message:       msg,
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
	})
	if err != nil {
		return false
	}
	return choice != "Yes"
}

// CheckUpdate compares this build against the latest GitHub release.
// Best-effort and fail-silent: the UI shows a link only when something
// newer exists. This is awareness, not an auto-updater.
func (a *App) CheckUpdate() update.Info {
	return update.Check(appVersion)
}

// DismissAlerts clears reviewed silent-install alerts.
func (a *App) DismissAlerts() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return config.Update(func(c *config.Config) error { c.Alerts = nil; return nil })
}

// RemovePackage uninstalls one package by name from a silent-install
// alert, and clears that alert once it is gone. Works unelevated for the
// current user.
func (a *App) RemovePackage(name string) error {
	if err := a.eng.Appx.Remove(name, elevate.IsElevated()); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return config.Update(func(c *config.Config) error { c.RemoveAlert(name); return nil })
}

// OpenLogFolder shows Deflater's log directory in Explorer.
func (a *App) OpenLogFolder() {
	cmd := exec.Command("explorer.exe", logging.LogDir())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
}

// installSelf copies the running exe to an admin-only per-machine
// directory so the auto-elevating scheduled task always runs a binary a
// non-admin cannot tamper with. Returns the stable path. Requires admin
// (its only caller, syncTask, is elevated). If a good copy already
// exists it is reused; a failed copy falls back to that existing copy
// before ever pointing the task at a volatile location.
func installSelf() (string, error) {
	src, err := os.Executable()
	if err != nil {
		return "", err
	}
	base := os.Getenv("ProgramFiles")
	if base == "" {
		base = `C:\Program Files`
	}
	dir := filepath.Join(base, "Deflater")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "Deflater.exe")
	if filepath.Clean(src) == filepath.Clean(dst) {
		return dst, nil
	}
	if err := copyFileAtomic(src, dst); err != nil {
		if _, statErr := os.Stat(dst); statErr == nil {
			logging.Logf("installSelf: copy failed, keeping existing %s: %v", dst, err)
			return dst, nil
		}
		return "", fmt.Errorf("install to %s failed: %w", dst, err)
	}
	return dst, nil
}

// copyFileAtomic writes src to a temp file beside dst and renames it in,
// so a partial copy never leaves a truncated exe the task would run.
func copyFileAtomic(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(dst), "Deflater-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}
