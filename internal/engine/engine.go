// Package engine turns catalog entries into action: reading each fix's
// live status, applying it, and undoing it.
package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"

	"deflater/internal/appx"
	"deflater/internal/catalog"
	"deflater/internal/logging"
	"deflater/internal/reg"
)

// Status values surfaced to the UI.
const (
	StatusOn        = "on"        // switch fully applied
	StatusOff       = "off"       // switch not applied
	StatusPartial   = "partial"   // some of a switch's values set
	StatusRemoved   = "removed"   // app not present
	StatusInstalled = "installed" // app present
)

type Engine struct {
	Appx *appx.Service
}

// Status reports the live state of a fix on this machine.
func (e *Engine) Status(f catalog.Fix) (string, error) {
	switch f.Kind {
	case catalog.Switch, catalog.OneDrive:
		set := 0
		for _, op := range f.Reg {
			ok, err := reg.IsSet(op)
			if err != nil {
				return "", err
			}
			if ok {
				set++
			}
		}
		switch {
		case set == len(f.Reg):
			return StatusOn, nil
		case set == 0:
			return StatusOff, nil
		default:
			return StatusPartial, nil
		}
	default: // app removals
		installed, err := e.Appx.Installed()
		if err != nil {
			return "", err
		}
		for _, pkg := range f.Appx {
			if installed[pkg] {
				return StatusInstalled, nil
			}
		}
		return StatusRemoved, nil
	}
}

// Capture records the pre-apply state of a fix's registry values so a
// later revert can restore exactly what the user had. Callers store this
// the first time a fix is applied and pass it back to Revert.
func (e *Engine) Capture(f catalog.Fix) ([]reg.Snapshot, error) {
	snaps := make([]reg.Snapshot, 0, len(f.Reg))
	for _, op := range f.Reg {
		s, err := reg.Capture(op)
		if err != nil {
			return nil, err
		}
		snaps = append(snaps, s)
	}
	return snaps, nil
}

// Apply puts the fix in effect. Elevated affects whether app removals
// can also be deprovisioned (blocked from returning for new accounts).
func (e *Engine) Apply(f catalog.Fix, elevated bool) error {
	for _, op := range f.Reg {
		if err := reg.Apply(op); err != nil {
			return err
		}
	}
	for _, pkg := range f.Appx {
		installed, err := e.Appx.Installed()
		if err != nil {
			return err
		}
		if !installed[pkg] {
			continue
		}
		if err := e.Appx.Remove(pkg, elevated); err != nil {
			return fmt.Errorf("remove %s: %w", pkg, err)
		}
	}
	if f.Kind == catalog.OneDrive {
		if err := uninstallOneDrive(); err != nil {
			return err
		}
	}
	logging.Logf("applied %s", f.ID)
	return nil
}

// Revert undoes a switch. When snapshots are provided (recorded at
// apply time) it restores the user's exact prior state; otherwise it
// falls back to the catalog's static revert. App removals cannot be
// undone here; the Store reinstalls them, and the UI says so.
func (e *Engine) Revert(f catalog.Fix, snapshots []reg.Snapshot) error {
	if f.Kind == catalog.AppJunk || f.Kind == catalog.AppMight {
		return fmt.Errorf("%s: app removals are reverted by reinstalling from the Microsoft Store", f.ID)
	}
	if len(snapshots) == len(f.Reg) && len(snapshots) > 0 {
		for _, s := range snapshots {
			if err := reg.Restore(s); err != nil {
				return err
			}
		}
		logging.Logf("reverted %s (restored snapshot)", f.ID)
		return nil
	}
	for _, op := range f.Reg {
		if err := reg.Undo(op); err != nil {
			return err
		}
	}
	logging.Logf("reverted %s (static default)", f.ID)
	return nil
}

// uninstallOneDrive stops the sync client and runs Microsoft's own
// uninstaller. Files in OneDrive's cloud are untouched; this removes
// only the local client. The uninstaller path is located from known
// system locations and the registered uninstall command, but is always
// validated to be a OneDriveSetup.exe that exists before running, and
// is launched with a normal argument (never a shell), so a tampered
// HKCU uninstall string cannot inject a command.
func uninstallOneDrive() error {
	kill := exec.Command("taskkill.exe", "/IM", "OneDrive.exe", "/F")
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = kill.Run() // fine if it was not running

	setup := findOneDriveSetup()
	if setup == "" {
		return nil // already gone
	}
	return runUninstaller(setup)
}

// findOneDriveSetup returns a trusted path to OneDriveSetup.exe, or "".
// It checks the fixed system locations first, then the executable named
// by the registered uninstall command, but only accepts a path whose
// base name is OneDriveSetup.exe and which exists on disk.
func findOneDriveSetup() string {
	candidates := []string{
		filepath.Join(os.Getenv("SystemRoot"), "System32", "OneDriveSetup.exe"),
		filepath.Join(os.Getenv("SystemRoot"), "SysWOW64", "OneDriveSetup.exe"),
	}
	if p := oneDriveSetupFromRegistry(); p != "" {
		candidates = append(candidates, p)
	}
	for _, c := range candidates {
		if !strings.EqualFold(filepath.Base(c), "OneDriveSetup.exe") {
			continue
		}
		if info, err := os.Stat(c); err == nil && !info.IsDir() {
			return c
		}
	}
	return ""
}

// oneDriveSetupFromRegistry extracts just the executable path from the
// per-user uninstall command Windows registered, discarding any
// arguments. The value is treated as untrusted input.
func oneDriveSetupFromRegistry() string {
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Uninstall\OneDriveSetup.exe`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	s, _, err := k.GetStringValue("UninstallString")
	if err != nil {
		return ""
	}
	s = strings.TrimSpace(s)
	// A quoted path: take what's inside the quotes.
	if strings.HasPrefix(s, `"`) {
		if end := strings.Index(s[1:], `"`); end >= 0 {
			return s[1 : 1+end]
		}
		return ""
	}
	// Unquoted: the exe path ends at ".exe"; drop any trailing args.
	if i := strings.Index(strings.ToLower(s), ".exe"); i >= 0 {
		return s[:i+4]
	}
	return s
}

// runUninstaller launches OneDriveSetup.exe /uninstall with normal argv
// (no shell), killing the setup process if it overruns.
func runUninstaller(setup string) error {
	cmd := exec.Command(setup, "/uninstall")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	if err := cmd.Start(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Minute):
		_ = cmd.Process.Kill()
		return fmt.Errorf("OneDrive uninstaller timed out")
	}
}
