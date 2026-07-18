// Package appx manages Microsoft Store apps through PowerShell, the
// supported management surface for them. Removals are per-user plus,
// when elevated, deprovisioning so Windows Update does not re-seed the
// app for new accounts. Everything removed here can be reinstalled from
// the Microsoft Store.
package appx

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"deflater/internal/psrun"
)

// Service caches the installed package list because enumerating it costs
// a second or two and the UI asks for many statuses at once.
type Service struct {
	mu        sync.Mutex
	installed map[string]bool
}

// Installed returns the set of installed (non-framework) package names,
// loading it on first use.
func (s *Service) Installed() (map[string]bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.installed == nil {
		if err := s.refreshLocked(); err != nil {
			return nil, err
		}
	}
	return s.installed, nil
}

// Refresh re-reads the installed package list.
func (s *Service) Refresh() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.refreshLocked()
}

func (s *Service) refreshLocked() error {
	out, err := psrun.Run(
		`Get-AppxPackage | Where-Object { -not $_.IsFramework } | Select-Object -ExpandProperty Name | ConvertTo-Json -Compress`,
		90*time.Second)
	if err != nil {
		return err
	}
	names, err := parseNames(out)
	if err != nil {
		return err
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	s.installed = set
	return nil
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
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	s.installed = set
}

// Remove uninstalls the package and, when elevated, removes it for all
// users and deprovisions it so it will not reappear for new accounts or
// after feature updates.
func (s *Service) Remove(name string, elevated bool) error {
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
	s.mu.Lock()
	if s.installed != nil {
		delete(s.installed, name)
	}
	s.mu.Unlock()
	return nil
}
