// Package engine turns catalog entries into action: reading each fix's
// live status, applying it, and undoing it.
package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// Revert undoes a switch. App removals cannot be undone here; the Store
// reinstalls them, and the UI says so.
func (e *Engine) Revert(f catalog.Fix) error {
	if f.Kind == catalog.AppJunk || f.Kind == catalog.AppMight {
		return fmt.Errorf("%s: app removals are reverted by reinstalling from the Microsoft Store", f.ID)
	}
	for _, op := range f.Reg {
		if err := reg.Undo(op); err != nil {
			return err
		}
	}
	logging.Logf("reverted %s", f.ID)
	return nil
}

// uninstallOneDrive stops the sync client and runs Microsoft's own
// uninstaller. Files in OneDrive's cloud are untouched; this removes
// only the local client. On 24H2 OneDrive installs per-user, so the
// user's own uninstall command is preferred, with the classic machine
// installer as fallback.
func uninstallOneDrive() error {
	kill := exec.Command("taskkill.exe", "/IM", "OneDrive.exe", "/F")
	kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	_ = kill.Run() // fine if it was not running

	if cmdLine := oneDriveUninstallString(); cmdLine != "" {
		return runCommandLine(cmdLine)
	}
	for _, c := range []string{
		filepath.Join(os.Getenv("SystemRoot"), "System32", "OneDriveSetup.exe"),
		filepath.Join(os.Getenv("SystemRoot"), "SysWOW64", "OneDriveSetup.exe"),
	} {
		if _, err := os.Stat(c); err == nil {
			return runCommandLine(`"` + c + `" /uninstall`)
		}
	}
	return nil // already gone
}

// oneDriveUninstallString reads the per-user uninstall command Windows
// itself registered for OneDrive, if present.
func oneDriveUninstallString() string {
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
	return s
}

func runCommandLine(cmdLine string) error {
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
		CmdLine:       `cmd.exe /C "` + cmdLine + `"`,
	}
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
