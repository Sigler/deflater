package reg

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

// Tests exercise the real registry mechanism against a throwaway key in
// the current user's hive; nothing outside it is touched.
const testPath = `Software\DeflaterTest\reg`

func cleanup(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		_ = registry.DeleteKey(registry.CURRENT_USER, testPath)
		_ = registry.DeleteKey(registry.CURRENT_USER, `Software\DeflaterTest`)
	})
}

func TestApplyIsSetUndoDelete(t *testing.T) {
	cleanup(t)
	op := Op{Hive: "HKCU", Path: testPath, Name: "DeleteMe", Value: 1, Revert: "delete"}

	if set, _ := IsSet(op); set {
		t.Fatal("value should not be set before Apply")
	}
	if err := Apply(op); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if set, _ := IsSet(op); !set {
		t.Fatal("value should be set after Apply")
	}
	if err := Undo(op); err != nil {
		t.Fatalf("Undo: %v", err)
	}
	if set, _ := IsSet(op); set {
		t.Fatal("value should be gone after Undo")
	}
}

func TestUndoRestoresExplicitDefault(t *testing.T) {
	cleanup(t)
	op := Op{Hive: "HKCU", Path: testPath, Name: "Toggle", Value: 0, Revert: "set:1"}

	if err := Apply(op); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if err := Undo(op); err != nil {
		t.Fatalf("Undo: %v", err)
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, testPath, registry.QUERY_VALUE)
	if err != nil {
		t.Fatalf("open key: %v", err)
	}
	defer k.Close()
	v, _, err := k.GetIntegerValue("Toggle")
	if err != nil {
		t.Fatalf("value should exist after set-revert: %v", err)
	}
	if v != 1 {
		t.Fatalf("revert should restore 1, got %d", v)
	}
}

func TestIsSetDistinguishesValues(t *testing.T) {
	cleanup(t)
	op := Op{Hive: "HKCU", Path: testPath, Name: "N", Value: 1, Revert: "delete"}
	other := Op{Hive: "HKCU", Path: testPath, Name: "N", Value: 2, Revert: "delete"}

	if err := Apply(other); err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if set, _ := IsSet(op); set {
		t.Fatal("IsSet must compare the value, not just existence")
	}
}

func TestUndoMissingValueIsFine(t *testing.T) {
	cleanup(t)
	op := Op{Hive: "HKCU", Path: testPath + `\absent`, Name: "Nope", Value: 1, Revert: "delete"}
	if err := Undo(op); err != nil {
		t.Fatalf("undoing something never applied must not error: %v", err)
	}
}

func TestUnknownHiveRejected(t *testing.T) {
	op := Op{Hive: "HKLM_TYPO", Path: testPath, Name: "X", Value: 1, Revert: "delete"}
	if err := Apply(op); err == nil {
		t.Fatal("unknown hive must be rejected")
	}
}
