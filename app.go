package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync/atomic"
	"syscall"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"deflater/internal/appx"
	"deflater/internal/catalog"
	"deflater/internal/config"
	"deflater/internal/elevate"
	"deflater/internal/engine"
	"deflater/internal/logging"
	"deflater/internal/schtask"
	"deflater/internal/toast"
)

// App is the bridge the frontend calls. Every exported method becomes a
// typed TypeScript function in frontend/wailsjs.
type App struct {
	ctx context.Context
	eng *engine.Engine
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
	// Pending is set when an apply was interrupted by an elevation
	// relaunch; the UI resumes it automatically.
	Pending *config.Pending `json:"pending"`
}

// GetReport gathers catalog, statuses, and settings. Statuses involve a
// Store package enumeration, so the first call takes a moment.
func (a *App) GetReport() (Report, error) {
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
}

// ApplyOutcome is the result of an Apply call. If NeedsElevation is set,
// nothing was changed; the request is saved and the UI should call
// ElevateNow to relaunch with admin rights and resume.
type ApplyOutcome struct {
	NeedsElevation bool        `json:"needsElevation"`
	Results        []FixResult `json:"results"`
}

// Apply enables and disables fixes. Enable ids are applied, disable ids
// reverted. Progress events stream to the UI as "apply:progress".
// Unelevated it changes nothing and asks for an elevated relaunch; the
// request is only saved once the user confirms (SaveAndElevate).
func (a *App) Apply(enable []string, disable []string) (ApplyOutcome, error) {
	cfg := config.Load()
	if !elevate.IsElevated() {
		return ApplyOutcome{NeedsElevation: true}, nil
	}

	var results []FixResult
	step := func(id string, do func(catalog.Fix) error, managedAfter bool) {
		fix, ok := catalog.ByID(id)
		if !ok {
			return
		}
		res := FixResult{ID: id, OK: true}
		if err := do(fix); err != nil {
			res.OK = false
			res.Error = err.Error()
			logging.Logf("apply: %s failed: %v", id, err)
		} else {
			cfg.SetManaged(id, managedAfter)
		}
		res.Status, _ = a.eng.Status(fix)
		results = append(results, res)
		if a.ctx != nil {
			wruntime.EventsEmit(a.ctx, "apply:progress", res)
		}
	}
	for _, id := range enable {
		step(id, func(f catalog.Fix) error { return a.eng.Apply(f, true) }, true)
	}
	for _, id := range disable {
		step(id, func(f catalog.Fix) error { return a.eng.Revert(f) }, false)
	}

	cfg.Pending = nil
	// Keep the maintenance task in sync while we hold admin rights.
	a.syncMaintenanceTask(&cfg)
	if err := config.Save(cfg); err != nil {
		return ApplyOutcome{}, err
	}
	return ApplyOutcome{Results: results}, nil
}

// SaveAndElevate stores the confirmed apply request, relaunches Deflater
// with admin rights (standard Windows prompt), and closes this instance.
// The elevated instance resumes the request automatically.
func (a *App) SaveAndElevate(enable []string, disable []string) error {
	cfg := config.Load()
	cfg.Pending = &config.Pending{Enable: enable, Disable: disable}
	if err := config.Save(cfg); err != nil {
		return err
	}
	if err := elevate.Relaunch(); err != nil {
		// UAC declined; forget the request so it cannot fire later.
		cfg.Pending = nil
		_ = config.Save(cfg)
		logging.Logf("elevate: relaunch failed or was declined: %v", err)
		return fmt.Errorf("Windows did not grant administrator rights")
	}
	wruntime.Quit(a.ctx)
	return nil
}

// SetMaintenance turns the scheduled maintenance task on or off. The
// preference always saves; registering the task itself needs admin
// rights, so unelevated it takes effect on the next elevated apply.
func (a *App) SetMaintenance(on bool) (bool, error) {
	cfg := config.Load()
	cfg.Maintenance = on
	applied := false
	if elevate.IsElevated() {
		a.syncMaintenanceTask(&cfg)
		applied = cfg.Maintenance == on // sync flips it back on failure
	}
	if err := config.Save(cfg); err != nil {
		return applied, err
	}
	return applied, nil
}

// syncMaintenanceTask makes the real scheduled task match cfg.Maintenance.
// Must be called elevated. On enable it first copies the exe to a stable
// per-machine location so the task survives the original download being
// moved or deleted.
func (a *App) syncMaintenanceTask(cfg *config.Config) {
	if cfg.Maintenance {
		exePath, err := installSelf()
		if err == nil {
			err = schtask.Install(exePath)
		}
		if err != nil {
			logging.Logf("maintenance: enable failed: %v", err)
			cfg.Maintenance = false
			return
		}
		logging.Logf("maintenance: task registered at %s", exePath)
		return
	}
	if err := schtask.Uninstall(); err != nil {
		logging.Logf("maintenance: disable failed: %v", err)
	}
}

// SetDirty records how many staged changes are awaiting apply.
func (a *App) SetDirty(n int) {
	a.dirty.Store(int32(n))
}

// beforeClose warns when the window is closing with staged changes.
// Returning true prevents the close.
func (a *App) beforeClose(ctx context.Context) bool {
	n := int(a.dirty.Load())
	if n <= 0 {
		return false
	}
	// This copy lives in Go because the dialog is native and can outlive
	// the webview; keep it in sync with the frontend string catalog.
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

// SetWatcher toggles silent-install alerts (used by maintenance runs).
func (a *App) SetWatcher(on bool) error {
	cfg := config.Load()
	cfg.WatcherEnabled = on
	return config.Save(cfg)
}

// DismissAlerts clears reviewed silent-install alerts.
func (a *App) DismissAlerts() error {
	cfg := config.Load()
	cfg.Alerts = nil
	return config.Save(cfg)
}

// RemovePackage uninstalls one package by name, used from a
// silent-install alert. Works unelevated for the current user.
func (a *App) RemovePackage(name string) error {
	return a.eng.Appx.Remove(name, elevate.IsElevated())
}

// OpenLogFolder shows Deflater's log directory in Explorer.
func (a *App) OpenLogFolder() {
	cmd := exec.Command("explorer.exe", logging.LogDir())
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
}

// installSelf copies the running exe to %LOCALAPPDATA%\Deflater\bin so
// the scheduled task has a stable target. Running from there already is
// fine; the copy is skipped.
func installSelf() (string, error) {
	src, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(logging.Dir(), "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, "Deflater.exe")
	if filepath.Clean(src) == filepath.Clean(dst) {
		return dst, nil
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(dst, data, 0o755); err != nil {
		// The copy can be locked by a previous instance; the original
		// location still works as a fallback target.
		logging.Logf("installSelf: copy failed, task will point at %s: %v", src, err)
		return src, nil
	}
	return dst, nil
}
