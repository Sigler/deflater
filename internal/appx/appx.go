// Package appx manages Microsoft Store apps through PowerShell, the
// supported management surface for them. Removals are per-user plus,
// when elevated, deprovisioning so Windows Update does not re-seed the
// app for new accounts. Everything removed here can be reinstalled from
// the Microsoft Store.
package appx

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"deflater/internal/psrun"
)

// packageName bounds what may be interpolated into a PowerShell command:
// the Appx identity charset. Anything else is rejected before it can
// reach the shell.
var packageName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9.\-_]*$`)

// Service caches the installed package list because enumerating it costs
// a second or two and the UI asks for many statuses at once. The cache
// map is treated as immutable once built: mutations swap in a new map,
// so a caller holding a returned reference always sees a stable snapshot
// (no concurrent read/write panics).
type Service struct {
	mu        sync.Mutex
	installed map[string]bool
	attempted bool  // enumeration has been tried (success or failure)
	loadErr   error // remembered so a slow failure is not retried per fix
}

// Installed returns the set of installed (non-framework) package names,
// enumerating once and then serving from cache. A failed enumeration is
// remembered and returned immediately rather than retried on every call,
// so a slow PowerShell failure cannot stack dozens of 90s timeouts. The
// returned map must not be mutated by the caller.
func (s *Service) Installed() (map[string]bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.attempted {
		s.refreshLocked()
	}
	return s.installed, s.loadErr
}

// Refresh forces a fresh enumeration on the next Installed call.
func (s *Service) Refresh() error {
	s.mu.Lock()
	s.attempted = false
	s.mu.Unlock()
	_, err := s.Installed()
	return err
}

func (s *Service) refreshLocked() {
	s.attempted = true
	s.installed = map[string]bool{}
	s.loadErr = nil
	out, err := psrun.Run(
		`Get-AppxPackage | Where-Object { -not $_.IsFramework } | Select-Object -ExpandProperty Name | ConvertTo-Json -Compress`,
		90*time.Second)
	if err != nil {
		s.loadErr = err
		return
	}
	names, err := parseNames(out)
	if err != nil {
		s.loadErr = err
		return
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	s.installed = set
}

// parseNames parses ConvertTo-Json output for the package name list: a
// JSON array normally, a bare JSON string when exactly one package
// matched, and empty output when none did.
func parseNames(out string) ([]string, error) {
	if out == "" {
		return nil, nil
	}
	if strings.HasPrefix(out, "[") {
		var names []string
		if err := json.Unmarshal([]byte(out), &names); err != nil {
			return nil, fmt.Errorf("parse package list: %w", err)
		}
		return names, nil
	}
	var one string
	if err := json.Unmarshal([]byte(out), &one); err != nil {
		return nil, fmt.Errorf("parse package list: %w", err)
	}
	return []string{one}, nil
}

// Prime seeds the installed-package cache without querying the system.
// Tests use it to control what "installed" means.
func (s *Service) Prime(names []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempted = true
	s.loadErr = nil
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	s.installed = set
}

// Remove uninstalls the package and, when elevated, removes it for all
// users and deprovisions it so it will not reappear for new accounts or
// after feature updates. The name is validated against the Appx
// identity charset before use, so it can never break out of the quoted
// PowerShell literal.
func (s *Service) Remove(name string, elevated bool) error {
	if !packageName.MatchString(name) {
		return fmt.Errorf("refusing to remove package with unexpected name %q", name)
	}
	script := fmt.Sprintf(
		`$ErrorActionPreference='Stop'; Get-AppxPackage -Name '%[1]s' | Remove-AppxPackage`, name)
	if elevated {
		script = fmt.Sprintf(
			`$ErrorActionPreference='Stop'; Get-AppxPackage -AllUsers -Name '%[1]s' | Remove-AppxPackage -AllUsers`+
				`; Get-AppxProvisionedPackage -Online | Where-Object DisplayName -eq '%[1]s' | ForEach-Object { Remove-AppxProvisionedPackage -Online -PackageName $_.PackageName -ErrorAction SilentlyContinue } | Out-Null`, name)
	}
	if _, err := psrun.Run(script, 3*time.Minute); err != nil {
		return err
	}
	// Copy-on-write: build a new map without the package and swap it in,
	// so any goroutine still reading the old map sees a stable snapshot.
	s.mu.Lock()
	if s.installed != nil {
		next := make(map[string]bool, len(s.installed))
		for k, v := range s.installed {
			if k != name {
				next[k] = v
			}
		}
		s.installed = next
	}
	s.mu.Unlock()
	return nil
}
