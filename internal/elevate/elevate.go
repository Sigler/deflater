// Package elevate answers "are we admin?" and relaunches Deflater with
// administrator rights when the user approves an apply that needs them.
package elevate

import (
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

// IsElevated reports whether the current process runs with admin rights.
func IsElevated() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}

// Relaunch starts this same executable elevated (triggering the standard
// Windows UAC prompt) with the given arguments.
func Relaunch(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	verb, _ := windows.UTF16PtrFromString("runas")
	file, _ := windows.UTF16PtrFromString(exe)
	argp, _ := windows.UTF16PtrFromString(strings.Join(args, " "))
	dirp, _ := windows.UTF16PtrFromString(cwd)
	return windows.ShellExecute(0, verb, file, argp, dirp, windows.SW_SHOWNORMAL)
}
