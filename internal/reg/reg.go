// Package reg reads and writes the registry values behind Deflater's
// switches. Every operation is a plain DWORD set, and every revert either
// deletes the value (restoring the Windows default) or writes an explicit
// number. Nothing else is ever touched.
package reg

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Op is one registry value a fix manages.
type Op struct {
	Hive  string `json:"hive"` // "HKLM" or "HKCU"
	Path  string `json:"path"`
	Name  string `json:"name"`
	Value uint32 `json:"value"`
	// Revert describes how to undo the op: "delete" removes the value so
	// Windows falls back to its default, "set:N" writes N instead.
	Revert string `json:"revert"`
}

func root(hive string) (registry.Key, error) {
	switch hive {
	case "HKLM":
		return registry.LOCAL_MACHINE, nil
	case "HKCU":
		return registry.CURRENT_USER, nil
	}
	return 0, fmt.Errorf("unknown hive %q", hive)
}

// Snapshot records a value's state before Deflater touched it, so a
// revert can restore exactly what the user had rather than an assumed
// Windows default. This matters most for consent-gated privacy values
// that vary by machine and region.
type Snapshot struct {
	Op      Op     `json:"op"`
	Existed bool   `json:"existed"`
	Prior   uint32 `json:"prior"`
}

// read returns the value's current data and whether it exists. A missing
// key or value is (0, false, nil); only real read failures error.
func read(op Op) (uint32, bool, error) {
	r, err := root(op.Hive)
	if err != nil {
		return 0, false, err
	}
	k, err := registry.OpenKey(r, op.Path, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return 0, false, nil
		}
		return 0, false, err
	}
	defer k.Close()
	v, _, err := k.GetIntegerValue(op.Name)
	if err == registry.ErrNotExist {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return uint32(v), true, nil
}

// IsSet reports whether the value currently matches what Apply would write.
func IsSet(op Op) (bool, error) {
	v, ok, err := read(op)
	if err != nil {
		return false, err
	}
	return ok && v == op.Value, nil
}

// Capture records the current state of the value for later restore.
func Capture(op Op) (Snapshot, error) {
	v, ok, err := read(op)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Op: op, Existed: ok, Prior: v}, nil
}

// Restore returns the value to its captured state: rewrite the prior
// value if it existed, otherwise delete it.
func Restore(s Snapshot) error {
	if s.Existed {
		return Apply(Op{Hive: s.Op.Hive, Path: s.Op.Path, Name: s.Op.Name, Value: s.Prior})
	}
	return deleteValue(s.Op)
}

// Apply writes the value, creating the key path if needed.
func Apply(op Op) error {
	r, err := root(op.Hive)
	if err != nil {
		return err
	}
	k, _, err := registry.CreateKey(r, op.Path, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open %s\\%s: %w", op.Hive, op.Path, err)
	}
	defer k.Close()
	if err := k.SetDWordValue(op.Name, op.Value); err != nil {
		return fmt.Errorf("set %s\\%s\\%s: %w", op.Hive, op.Path, op.Name, err)
	}
	return nil
}

// Undo reverses the op according to its Revert mode. Prefer restoring a
// captured Snapshot (engine does when one exists); this static revert is
// the fallback for fixes applied before snapshots were recorded.
func Undo(op Op) error {
	if after, ok := strings.CutPrefix(op.Revert, "set:"); ok {
		n, err := strconv.ParseUint(after, 10, 32)
		if err != nil {
			return fmt.Errorf("bad revert spec %q: %w", op.Revert, err)
		}
		return Apply(Op{Hive: op.Hive, Path: op.Path, Name: op.Name, Value: uint32(n)})
	}
	// Default: delete the value, Windows falls back to its own default.
	return deleteValue(op)
}

// deleteValue removes the value. A missing key or value is success;
// every other failure (notably access-denied when not elevated) is a
// real error, so the UI never reports a switch reverted when it wasn't.
func deleteValue(op Op) error {
	r, err := root(op.Hive)
	if err != nil {
		return err
	}
	k, err := registry.OpenKey(r, op.Path, registry.SET_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return fmt.Errorf("open %s\\%s: %w", op.Hive, op.Path, err)
	}
	defer k.Close()
	if err := k.DeleteValue(op.Name); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("delete %s\\%s\\%s: %w", op.Hive, op.Path, op.Name, err)
	}
	return nil
}
