// Package maintain is the headless mode behind "Deflater.exe
// --maintenance", run by the scheduled task after sign-in and weekly.
// It re-applies drifted fixes (Windows updates love to re-seed junk) and
// watches for apps that appeared without the user asking.
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

// Run performs one maintenance pass. It returns the number of fixes it
// had to re-apply.
func Run() int {
	logging.Logf("maintenance: start (elevated=%v)", elevate.IsElevated())
	cfg := config.Load()
	eng := &engine.Engine{Appx: &appx.Service{}}
	elevated := elevate.IsElevated()

	reapplied := 0
	for _, id := range cfg.Managed {
		fix, ok := catalog.ByID(id)
		if !ok {
			continue // fix retired in a newer version
		}
		status, err := eng.Status(fix)
		if err != nil {
			logging.Logf("maintenance: %s: status check failed: %v", id, err)
			continue
		}
		needsWork := status == engine.StatusOff || status == engine.StatusPartial || status == engine.StatusInstalled
		if !needsWork {
			continue
		}
		if err := eng.Apply(fix, elevated); err != nil {
			logging.Logf("maintenance: %s: re-apply failed: %v", id, err)
			continue
		}
		logging.Logf("maintenance: re-applied %s (was %s)", id, status)
		reapplied++
	}

	runWatcher(&cfg, eng)

	if err := config.Save(cfg); err != nil {
		logging.Logf("maintenance: save config failed: %v", err)
	}
	logging.Logf("maintenance: done, re-applied %d", reapplied)
	return reapplied
}

// runWatcher diffs the installed apps against the last snapshot and
// records alerts (plus a toast) for anything that arrived on its own.
func runWatcher(cfg *config.Config, eng *engine.Engine) {
	current, err := eng.Appx.Installed()
	if err != nil {
		logging.Logf("watcher: package list failed: %v", err)
		return
	}
	if cfg.WatcherEnabled {
		arrivals := watcher.NewArrivals(current, cfg.Snapshot, catalog.ManagedPackages(cfg.Managed))
		for _, pkg := range arrivals {
			logging.Logf("watcher: new app appeared: %s", pkg)
			cfg.AddAlert(pkg)
		}
		if len(arrivals) > 0 {
			body := arrivals[0]
			if len(arrivals) > 1 {
				body = fmt.Sprintf("%s and %d more", arrivals[0], len(arrivals)-1)
			}
			if err := toast.Show("An app installed itself", body+" appeared without you asking. Open Deflater to review or remove it."); err != nil {
				logging.Logf("watcher: toast failed: %v", err)
			}
		}
	}
	cfg.Snapshot = watcher.SnapshotOf(current)
}
