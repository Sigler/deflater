package config

import (
	"os"

	"golang.org/x/sys/windows"
)

// acquireLock takes an exclusive advisory lock on a dedicated lock file,
// blocking until it is available, and returns a function that releases
// it. It coordinates the GUI and the maintenance process so their
// read-modify-write cycles on config.json cannot interleave.
func acquireLock() (func(), error) {
	f, err := os.OpenFile(lockPath(), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	h := windows.Handle(f.Fd())
	var overlapped windows.Overlapped
	// LOCKFILE_EXCLUSIVE_LOCK, blocking (no _FAIL_IMMEDIATELY), over the
	// whole file.
	err = windows.LockFileEx(h, windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &overlapped)
	if err != nil {
		f.Close()
		return nil, err
	}
	return func() {
		_ = windows.UnlockFileEx(h, 0, 1, 0, &overlapped)
		f.Close()
	}, nil
}
