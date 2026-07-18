// Package psrun executes PowerShell snippets in a hidden window.
// Deflater leans on PowerShell for the operations Windows only exposes
// well through it: Store app management, scheduled tasks, and toasts.
package psrun

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const createNoWindow = 0x08000000

// Run executes a PowerShell script and returns its combined output.
// The window is hidden and the call is killed after timeout.
func Run(script string, timeout time.Duration) (string, error) {
	cmd := exec.Command("powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: createNoWindow}

	done := make(chan struct{})
	var out []byte
	var err error
	go func() {
		out, err = cmd.CombinedOutput()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		return "", fmt.Errorf("powershell timed out after %s", timeout)
	}

	text := strings.TrimSpace(string(out))
	if err != nil {
		return text, fmt.Errorf("powershell failed: %w (output: %s)", err, text)
	}
	return text, nil
}
