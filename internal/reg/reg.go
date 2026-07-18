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

// IsSet reports whether the value currently matches what Apply would write.
func IsSet(op Op) (bool, error) {
	r, err := root(op.Hive)
	if err != nil {
		return false, err
	}
	k, err := registry.OpenKey(r, op.Path, registry.QUERY_VALUE)
	if err != nil {
		return false, nil // key absent means not applied, not an error
	}
	defer k.Close()
	v, _, err := k.GetIntegerValue(op.Name)
	if err != nil {
		return false, nil
	}
	return uint32(v) == op.Value, nil
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

// Undo reverses the op according to its Revert mode.
func Undo(op Op) error {
	r, err := root(op.Hive)
	if err != nil {
		return err
	}
	if after, ok := strings.CutPrefix(op.Revert, "set:"); ok {
		n, err := strconv.ParseUint(after, 10, 32)
		if err != nil {
			return fmt.Errorf("bad revert spec %q: %w", op.Revert, err)
		}
		k, _, err := registry.CreateKey(r, op.Path, registry.SET_VALUE)
		if err != nil {
			return err
		}
		defer k.Close()
		return k.SetDWordValue(op.Name, uint32(n))
	}
	// Default: delete the value, Windows falls back to its own default.
	k, err := registry.OpenKey(r, op.Path, registry.SET_VALUE)
	if err != nil {
		return nil // key already gone
	}
	defer k.Close()
	if err := k.DeleteValue(op.Name); err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("delete %s\\%s\\%s: %w", op.Hive, op.Path, op.Name, err)
	}
	return nil
}
