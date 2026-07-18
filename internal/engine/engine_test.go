package engine

import (
	"testing"

	"golang.org/x/sys/windows/registry"

	"deflater/internal/appx"
	"deflater/internal/catalog"
	"deflater/internal/reg"
)

// Registry tests run against a throwaway key in the current user's
// hive; the app-status tests use a primed package cache. Nothing else
// on the machine is touched.
const testPath = `Software\DeflaterTest\engine`

func cleanup(t *testing.T) {
	t.Helper()
	t.Setenv("DEFLATER_DATA_DIR", t.TempDir()) // keep logs out of the real data dir
	t.Cleanup(func() {
		_ = registry.DeleteKey(registry.CURRENT_USER, testPath)
		_ = registry.DeleteKey(registry.CURRENT_USER, `Software\DeflaterTest`)
	})
}

func switchFix() catalog.Fix {
	return catalog.Fix{
		ID:   "test-switch",
		Kind: catalog.Switch,
		Reg: []reg.Op{
			{Hive: "HKCU", Path: testPath, Name: "A", Value: 1, Revert: "delete"},
			{Hive: "HKCU", Path: testPath, Name: "B", Value: 0, Revert: "set:1"},
		},
	}
}

func TestSwitchLifecycle(t *testing.T) {
	cleanup(t)
	eng := &Engine{Appx: &appx.Service{}}
	fix := switchFix()

	status, err := eng.Status(fix)
	if err != nil || status != StatusOff {
		t.Fatalf("before apply: status %q err %v", status, err)
	}

	if err := eng.Apply(fix, false); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if status, _ = eng.Status(fix); status != StatusOn {
		t.Fatalf("after apply: status %q", status)
	}

	// Simulate drift: one of the two values reset by something else.
	k, err := registry.OpenKey(registry.CURRENT_USER, testPath, registry.SET_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	if err := k.DeleteValue("A"); err != nil {
		t.Fatal(err)
	}
	k.Close()
	if status, _ = eng.Status(fix); status != StatusPartial {
		t.Fatalf("after drift: status %q, want partial", status)
	}

	if err := eng.Revert(fix, nil); err != nil {
		t.Fatalf("Revert: %v", err)
	}
	if status, _ = eng.Status(fix); status != StatusOff {
		t.Fatalf("after revert: status %q", status)
	}

	// The set:1 revert must have restored B to its Windows default.
	k, err = registry.OpenKey(registry.CURRENT_USER, testPath, registry.QUERY_VALUE)
	if err != nil {
		t.Fatal(err)
	}
	defer k.Close()
	if v, _, err := k.GetIntegerValue("B"); err != nil || v != 1 {
		t.Fatalf("B after revert = %d (err %v), want 1", v, err)
	}
}

func TestRevertRestoresCapturedSnapshot(t *testing.T) {
	cleanup(t)
	eng := &Engine{Appx: &appx.Service{}}
	op := reg.Op{Hive: "HKCU", Path: testPath, Name: "Consent", Value: 0, Revert: "set:1"}
	fix := catalog.Fix{ID: "test-consent", Kind: catalog.Switch, Reg: []reg.Op{op}}

	// The user's real prior state: the value did not exist at all.
	snaps, err := eng.Capture(fix)
	if err != nil {
		t.Fatalf("Capture: %v", err)
	}
	if err := eng.Apply(fix, false); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	// Reverting with the snapshot must delete the value (restore absence),
	// NOT write the static default of 1, which would flip a consent
	// setting on that the user never had.
	if err := eng.Revert(fix, snaps); err != nil {
		t.Fatalf("Revert: %v", err)
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, testPath, registry.QUERY_VALUE)
	if err == nil {
		if _, _, e := k.GetIntegerValue("Consent"); e == nil {
			t.Fatal("snapshot revert should have deleted the value, but it exists")
		}
		k.Close()
	}
}

func TestAppFixStatus(t *testing.T) {
	cleanup(t)
	svc := &appx.Service{}
	eng := &Engine{Appx: svc}
	fix := catalog.Fix{ID: "test-app", Kind: catalog.AppJunk, Appx: []string{"Fake.App"}}

	svc.Prime([]string{"Fake.App", "Other.App"})
	if status, err := eng.Status(fix); err != nil || status != StatusInstalled {
		t.Fatalf("primed installed: status %q err %v", status, err)
	}

	svc.Prime(nil)
	if status, err := eng.Status(fix); err != nil || status != StatusRemoved {
		t.Fatalf("primed empty: status %q err %v", status, err)
	}
}

func TestRevertRefusesAppRemovals(t *testing.T) {
	cleanup(t)
	eng := &Engine{Appx: &appx.Service{}}
	fix := catalog.Fix{ID: "test-app", Kind: catalog.AppMight, Appx: []string{"Fake.App"}}
	if err := eng.Revert(fix, nil); err == nil {
		t.Fatal("reverting an app removal must error; the Store reinstalls apps")
	}
}
