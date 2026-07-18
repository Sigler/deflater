package maintain

import (
	"errors"
	"strings"
	"testing"

	"deflater/internal/catalog"
	"deflater/internal/config"
)

// fake wires the maintenance pass to in-memory state so the whole
// decision path runs without touching the machine.
type fake struct {
	fixes     map[string]catalog.Fix
	statuses  map[string]string
	applyErr  map[string]error
	installed map[string]bool
	applied   []string
	notified  []string
}

func (f *fake) deps() deps {
	return deps{
		byID: func(id string) (catalog.Fix, bool) {
			fx, ok := f.fixes[id]
			return fx, ok
		},
		status: func(fx catalog.Fix) (string, error) {
			s, ok := f.statuses[fx.ID]
			if !ok {
				return "", errors.New("status unavailable")
			}
			return s, nil
		},
		apply: func(fx catalog.Fix, _ bool) error {
			if err := f.applyErr[fx.ID]; err != nil {
				return err
			}
			f.applied = append(f.applied, fx.ID)
			return nil
		},
		installed: func() (map[string]bool, error) { return f.installed, nil },
		notify: func(_, body string) error {
			f.notified = append(f.notified, body)
			return nil
		},
	}
}

func quietLogs(t *testing.T) {
	t.Helper()
	t.Setenv("DEFLATER_DATA_DIR", t.TempDir())
}

func fix(id string) catalog.Fix { return catalog.Fix{ID: id, Kind: catalog.Switch} }

func TestReappliesOnlyDriftedManagedFixes(t *testing.T) {
	quietLogs(t)
	f := &fake{
		fixes:     map[string]catalog.Fix{"a": fix("a"), "b": fix("b"), "c": fix("c")},
		statuses:  map[string]string{"a": "on", "b": "off", "c": "partial"},
		installed: map[string]bool{},
	}
	cfg := &config.Config{Maintenance: true, Managed: []string{"a", "b", "c", "retired-id"}}

	n := run(cfg, f.deps())

	if n != 2 {
		t.Fatalf("reapplied %d, want 2", n)
	}
	if len(f.applied) != 2 || f.applied[0] != "b" || f.applied[1] != "c" {
		t.Fatalf("applied %v, want [b c]", f.applied)
	}
}

func TestWatcherOnlyModeNeverApplies(t *testing.T) {
	quietLogs(t)
	f := &fake{
		fixes:     map[string]catalog.Fix{"b": fix("b")},
		statuses:  map[string]string{"b": "off"},
		installed: map[string]bool{"Some.App": true},
	}
	cfg := &config.Config{Maintenance: false, WatcherEnabled: true, Managed: []string{"b"}}

	if n := run(cfg, f.deps()); n != 0 || len(f.applied) != 0 {
		t.Fatalf("watcher-only mode applied fixes: n=%d applied=%v", n, f.applied)
	}
	if len(cfg.Snapshot) != 1 {
		t.Fatalf("watcher-only mode must still snapshot, got %v", cfg.Snapshot)
	}
}

func TestWatcherAlertsAndNotifiesOnArrivals(t *testing.T) {
	quietLogs(t)
	f := &fake{
		fixes:     map[string]catalog.Fix{},
		installed: map[string]bool{"Old.App": true, "LG.ThinQ": true, "Another.Surprise": true},
	}
	cfg := &config.Config{WatcherEnabled: true, Snapshot: []string{"Old.App"}}

	run(cfg, f.deps())

	if len(cfg.Alerts) != 2 {
		t.Fatalf("alerts %v, want 2", cfg.Alerts)
	}
	if len(f.notified) != 1 || !strings.Contains(f.notified[0], "and 1 more") {
		t.Fatalf("notification %v, want one toast mentioning 'and 1 more'", f.notified)
	}
}

func TestWatcherIgnoresPackagesDeflaterManages(t *testing.T) {
	quietLogs(t)
	// A managed package coming back is maintenance's job to remove,
	// not the watcher's job to announce.
	f := &fake{
		fixes: map[string]catalog.Fix{
			"remove-news": {ID: "remove-news", Kind: catalog.AppJunk, Appx: []string{"Managed.Pkg"}},
		},
		statuses:  map[string]string{"remove-news": "removed"},
		installed: map[string]bool{"Managed.Pkg": true},
	}
	cfg := &config.Config{Maintenance: true, WatcherEnabled: true, Managed: []string{"remove-news"}, Snapshot: []string{}}

	run(cfg, f.deps())

	if len(cfg.Alerts) != 0 {
		t.Fatalf("managed package alerted: %v", cfg.Alerts)
	}
}

func TestFirstRunEstablishesBaselineWithoutAlerting(t *testing.T) {
	quietLogs(t)
	f := &fake{fixes: map[string]catalog.Fix{}, installed: map[string]bool{"Preexisting.App": true}}
	cfg := &config.Config{WatcherEnabled: true} // no snapshot yet

	run(cfg, f.deps())

	if len(cfg.Alerts) != 0 || len(f.notified) != 0 {
		t.Fatalf("first run must not alert: alerts=%v notified=%v", cfg.Alerts, f.notified)
	}
	if len(cfg.Snapshot) != 1 {
		t.Fatalf("first run must establish the baseline, got %v", cfg.Snapshot)
	}
}

func TestErrorsSkipFixButContinuePass(t *testing.T) {
	quietLogs(t)
	f := &fake{
		fixes:     map[string]catalog.Fix{"broken": fix("broken"), "erroring": fix("erroring"), "fine": fix("fine")},
		statuses:  map[string]string{"erroring": "off", "fine": "off"}, // "broken" has no status: status() errors
		applyErr:  map[string]error{"erroring": errors.New("apply exploded")},
		installed: map[string]bool{},
	}
	cfg := &config.Config{Maintenance: true, Managed: []string{"broken", "erroring", "fine"}}

	if n := run(cfg, f.deps()); n != 1 || len(f.applied) != 1 || f.applied[0] != "fine" {
		t.Fatalf("pass must survive per-fix failures: n=%d applied=%v", n, f.applied)
	}
}
