package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func useTempDir(t *testing.T) {
	t.Helper()
	t.Setenv("DEFLATER_DATA_DIR", t.TempDir())
}

func TestDefaultsWhenNoFile(t *testing.T) {
	useTempDir(t)
	c := Load()
	if !c.WatcherEnabled {
		t.Fatal("watcher should default on (the block-and-warn default)")
	}
	if c.Maintenance {
		t.Fatal("maintenance should default off")
	}
}

func TestRoundTrip(t *testing.T) {
	useTempDir(t)
	c := Config{
		Managed:        []string{"lockscreen-ads", "app-news"},
		Maintenance:    true,
		WatcherEnabled: false,
		Snapshot:       []string{"A.App"},
		Pending:        &Pending{Enable: []string{"x"}, Disable: []string{"y"}},
	}
	c.AddAlert("LG.ThinQ")
	if err := Save(c); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got := Load()
	if !reflect.DeepEqual(got.Managed, c.Managed) ||
		got.Maintenance != c.Maintenance ||
		got.WatcherEnabled != c.WatcherEnabled ||
		!reflect.DeepEqual(got.Snapshot, c.Snapshot) ||
		len(got.Alerts) != 1 || got.Alerts[0].Package != "LG.ThinQ" ||
		got.Pending == nil || got.Pending.Enable[0] != "x" {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}

func TestCorruptFileFallsBackToDefaults(t *testing.T) {
	useTempDir(t)
	if err := os.WriteFile(filepath.Join(os.Getenv("DEFLATER_DATA_DIR"), "config.json"), []byte("{nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	c := Load()
	if !c.WatcherEnabled || c.Maintenance || len(c.Managed) != 0 {
		t.Fatalf("corrupt config must yield defaults, got %+v", c)
	}
}

func TestAddAlertDeduplicates(t *testing.T) {
	c := Config{}
	c.AddAlert("LG.ThinQ")
	c.AddAlert("LG.ThinQ")
	if len(c.Alerts) != 1 {
		t.Fatalf("duplicate alerts recorded: %v", c.Alerts)
	}
}

func TestSetManaged(t *testing.T) {
	c := Config{}
	c.SetManaged("a", true)
	c.SetManaged("b", true)
	c.SetManaged("a", true) // no duplicate
	if !reflect.DeepEqual(c.Managed, []string{"b", "a"}) && !reflect.DeepEqual(c.Managed, []string{"a", "b"}) {
		if len(c.Managed) != 2 || !c.HasManaged("a") || !c.HasManaged("b") {
			t.Fatalf("managed after adds: %v", c.Managed)
		}
	}
	c.SetManaged("a", false)
	if c.HasManaged("a") || !c.HasManaged("b") {
		t.Fatalf("managed after remove: %v", c.Managed)
	}
}
