// Package config persists Deflater's state to
// %LOCALAPPDATA%\Deflater\config.json: which fixes the user manages,
// whether maintenance is on, and the app snapshot the watcher diffs
// against to catch silent installs.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"deflater/internal/logging"
)

// Alert records an app that appeared without the user asking for it.
type Alert struct {
	Package string `json:"package"`
	Seen    string `json:"seen"` // RFC 3339, kept as a string for clean bindings
}

// Pending carries an apply request across an elevation relaunch.
type Pending struct {
	Enable  []string `json:"enable"`
	Disable []string `json:"disable"`
}

type Config struct {
	// Managed lists fix ids the user has applied; maintenance re-applies
	// them when Windows drifts (for example after a feature update).
	Managed []string `json:"managed"`
	// Maintenance mirrors whether the scheduled task is meant to exist.
	Maintenance bool `json:"maintenance"`
	// WatcherEnabled controls the silent-install toast alerts.
	WatcherEnabled bool `json:"watcherEnabled"`
	// Snapshot is the installed-app list from the last watcher run.
	Snapshot []string `json:"snapshot,omitempty"`
	// Alerts are silent installs noticed but not yet reviewed in the app.
	Alerts []Alert `json:"alerts,omitempty"`
	// Pending is a saved apply request waiting for an elevated relaunch.
	Pending *Pending `json:"pending,omitempty"`
}

func path() string { return filepath.Join(logging.Dir(), "config.json") }

// Load reads the config, returning sensible defaults when absent.
func Load() Config {
	c := Config{WatcherEnabled: true}
	data, err := os.ReadFile(path())
	if err != nil {
		return c
	}
	if err := json.Unmarshal(data, &c); err != nil {
		logging.Logf("config: unreadable, starting fresh: %v", err)
		return Config{WatcherEnabled: true}
	}
	return c
}

// Save writes the config atomically.
func Save(c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	tmp := path() + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path())
}

// HasManaged reports whether id is in the managed list.
func (c *Config) HasManaged(id string) bool {
	for _, m := range c.Managed {
		if m == id {
			return true
		}
	}
	return false
}

// AddAlert records a silent install if it is not already recorded.
func (c *Config) AddAlert(pkg string) {
	for _, a := range c.Alerts {
		if a.Package == pkg {
			return
		}
	}
	c.Alerts = append(c.Alerts, Alert{Package: pkg, Seen: time.Now().Format(time.RFC3339)})
}

// SetManaged adds or removes id from the managed list.
func (c *Config) SetManaged(id string, on bool) {
	out := c.Managed[:0]
	for _, m := range c.Managed {
		if m != id {
			out = append(out, m)
		}
	}
	c.Managed = out
	if on {
		c.Managed = append(c.Managed, id)
	}
}
