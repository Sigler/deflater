// Package watcher spots apps that appear without the user asking:
// manufacturer auto-installs (the LG case), promo re-seeds after
// updates, and similar. It compares the current installed-app list
// against the snapshot from the previous maintenance run.
package watcher

import (
	"sort"
	"strings"
)

// System packages that come and go with normal Windows servicing.
// Changes to these are routine, not silent bloat, so they never alert.
var routinePrefixes = []string{
	"Microsoft.NET.",
	"Microsoft.UI.",
	"Microsoft.VCLibs",
	"Microsoft.WindowsAppRuntime",
	"Microsoft.DirectX",
	"Microsoft.Services.Store",
	"Microsoft.StorePurchaseApp",
	"Microsoft.WindowsStore",
	"Microsoft.SecHealthUI",
	"Microsoft.AAD.",
	"Microsoft.AccountsControl",
	"Microsoft.Win32WebViewHost",
	"Microsoft.Windows.",
	"MicrosoftWindows.Client",
	"MicrosoftWindows.UndockedDevKit",
	"windows.immersivecontrolpanel",
	"Microsoft.ApplicationCompatibilityEnhancements",
	"Microsoft.AV1VideoExtension",
	"Microsoft.AVCEncoderVideoExtension",
	"Microsoft.HEIFImageExtension",
	"Microsoft.HEVCVideoExtension",
	"Microsoft.RawImageExtension",
	"Microsoft.VP9VideoExtensions",
	"Microsoft.WebMediaExtensions",
	"Microsoft.WebpImageExtension",
}

func routine(name string) bool {
	for _, p := range routinePrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// NewArrivals returns packages present now that were not in the snapshot,
// excluding routine system packages and anything Deflater itself manages.
// A nil snapshot means "first run": nothing to compare against yet.
func NewArrivals(current map[string]bool, snapshot []string, managed map[string]bool) []string {
	if snapshot == nil {
		return nil
	}
	prev := make(map[string]bool, len(snapshot))
	for _, n := range snapshot {
		prev[n] = true
	}
	var out []string
	for name := range current {
		if !prev[name] && !routine(name) && !managed[name] {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// SnapshotOf flattens the installed set into a sorted list for storage.
func SnapshotOf(current map[string]bool) []string {
	out := make([]string, 0, len(current))
	for name := range current {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
