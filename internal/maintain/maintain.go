// Package maintain is the headless mode behind "Deflater.exe
// --maintenance", run by the scheduled task after sign-in and weekly.
// It re-applies drifted fixes (Windows updates love to re-seed junk) and
// watches for apps that appear without the user asking.
package maintain

import (
	"fmt"

	"deflater/internal/appx"
	"deflater/internal/catalog"
	"deflater/internal/config"
	"deflater/internal/elevate"
	"deflater/internal/engine"
	"deflater/internal/logging"
	"deflater/internal/toast"
	"deflater/internal/watcher"
)

// deps are the seams between the maintenance logic and the real system.
// Run wires the real implementations; tests inject fakes, so the whole
// decision path runs under test without touching the machine.
type deps struct {
	byID      func(id string) (catalog.Fix, bool)
	status    func(f catalog.Fix) (string, error)
	apply     func(f catalog.Fix, elevated bool) error
	installed func() (map[string]bool, error)
	notify    func(title, body string) error
	elevated  bool
}

// Run performs one maintenance pass. It returns the number of fixes it
// had to re-apply.
func Run() int {
	eng := &engine.Engine{Appx: &appx.Service{}}
	d := deps{
		byID:      catalog.ByID,
		status:    eng.Status,
		apply:     eng.Apply,
		installed: eng.Appx.Installed,
		notify:    toast.Show,
		elevated:  elevate.IsElevated(),
	}
	cfg := config.Load()
	logging.Logf("maintenance: start (elevated=%v, maintenance=%v, watcher=%v)",
		d.elevated, cfg.Maintenance, cfg.WatcherEnabled)
	before := len(cfg.Alerts)
	n := run(&cfg, d)

	// The long work above ran against a loaded copy. Persist only the two
	// fields maintenance owns — the fresh snapshot and any newly-found
	// alerts — under the lock, so a GUI change made meanwhile survives.
	newAlerts := append([]config.Alert(nil), cfg.Alerts[before:]...)
	snapshot := cfg.Snapshot
	if err := config.Update(func(c *config.Config) error {
		c.Snapshot = snapshot
		for _, al := range newAlerts {
			c.AddAlert(al.Package)
		}
		return nil
	}); err != nil {
		logging.Logf("maintenance: save config failed: %v", err)
	}
	logging.Logf("maintenance: done, re-applied %d", n)
	return n
}

// run is the whole pass: re-apply drifted managed fixes (when
// maintenance is on), then check for silent installs (when the watcher
// is on). The scheduled task exists if either switch wants it.
func run(cfg *config.Config, d deps) int {
	reapplied := 0
	if cfg.Maintenance {
		for _, id := range cfg.Managed {
			fix, ok := d.byID(id)
			if !ok {
				continue // fix retired in a newer version
			}
			status, err := d.status(fix)
			if err != nil {
				logging.Logf("maintenance: %s: status check failed: %v", id, err)
				continue
			}
			if !needsReapply(status) {
				continue
			}
			if err := d.apply(fix, d.elevated); err != nil {
				logging.Logf("maintenance: %s: re-apply failed: %v", id, err)
				continue
			}
			logging.Logf("maintenance: re-applied %s (was %s)", id, status)
			reapplied++
		}
	}
	runWatcher(cfg, d)
	return reapplied
}

// needsReapply reports whether a managed fix has drifted out of effect.
func needsReapply(status string) bool {
	return status == engine.StatusOff || status == engine.StatusPartial || status == engine.StatusInstalled
}

// runWatcher diffs the installed apps against the last snapshot and
// records alerts (plus a notification) for anything that arrived on its
// own. The snapshot updates even when alerts are off, so turning the
// watcher on later starts from current reality.
func runWatcher(cfg *config.Config, d deps) {
	current, err := d.installed()
	if err != nil {
		logging.Logf("watcher: package list failed: %v", err)
		return
	}
	if cfg.WatcherEnabled {
		managedPkgs := map[string]bool{}
		for _, id := range cfg.Managed {
			if f, ok := d.byID(id); ok {
				for _, p := range f.Appx {
					managedPkgs[p] = true
				}
			}
		}
		arrivals := watcher.NewArrivals(current, cfg.Snapshot, managedPkgs)
		for _, pkg := range arrivals {
			logging.Logf("watcher: new app appeared: %s", pkg)
			cfg.AddAlert(pkg)
		}
		if len(arrivals) > 0 {
			body := arrivals[0]
			if len(arrivals) > 1 {
				body = fmt.Sprintf("%s and %d more", arrivals[0], len(arrivals)-1)
			}
			if err := d.notify("An app installed itself", body+" appeared without you asking. Open Deflater to review or remove it."); err != nil {
				logging.Logf("watcher: notification failed: %v", err)
			}
		}
	}
	cfg.Snapshot = watcher.SnapshotOf(current)
}
