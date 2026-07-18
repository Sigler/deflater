// Package config persists Deflater's state to
// %LOCALAPPDATA%\Deflater\config.json: which fixes the user manages,
// whether maintenance is on, the app snapshot the watcher diffs against
// to catch silent installs, and the registry values captured before each
// fix was applied so a revert can restore them exactly.
//
// Two processes touch this file: the GUI and the scheduled maintenance
// pass. Every read-modify-write goes through Update, which holds a
// cross-process file lock for the whole cycle so neither can clobber the
// other's changes.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"deflater/internal/logging"
	"deflater/internal/reg"
)

// Alert records an app that appeared without the user asking for it.
type Alert struct {
	Package string `json:"package"`
	Seen    string `json:"seen"` // RFC 3339, kept as a string for clean bindings
}

// Pending carries an apply request across an elevation relaunch. Token
// and Created let the elevated instance confirm it is resuming the
// request this session created, and reject a stale one.
type Pending struct {
	Enable  []string `json:"enable"`
	Disable []string `json:"disable"`
	Token   string   `json:"token"`
	Created string   `json:"created"` // RFC 3339
}

// Expired reports whether a pending request is older than the window in
// which an elevation handoff could plausibly still be in progress.
func (p *Pending) Expired(now time.Time) bool {
	t, err := time.Parse(time.RFC3339, p.Created)
	if err != nil {
		return true
	}
	return now.Sub(t) > 10*time.Minute
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
	// Snapshots holds the pre-apply registry state for each managed fix,
	// keyed by fix id, so revert restores the user's exact prior values.
	Snapshots map[string][]reg.Snapshot `json:"snapshots,omitempty"`
}

func path() string     { return filepath.Join(logging.Dir(), "config.json") }
func lockPath() string { return filepath.Join(logging.Dir(), "config.json.lock") }

func defaults() Config { return Config{WatcherEnabled: true} }

// Load reads the config, returning sensible defaults when absent. A
// corrupt file is preserved (renamed aside) rather than silently
// overwritten, so a transient bad read never erases the user's state.
func Load() Config {
	data, err := os.ReadFile(path())
	if err != nil {
		return defaults()
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		bad := path() + ".bad-" + time.Now().Format("20060102-150405")
		_ = os.Rename(path(), bad)
		logging.Logf("config: unreadable, preserved as %s, starting fresh: %v", bad, err)
		return defaults()
	}
	return c
}

// Save writes the config atomically: a uniquely-named temp file, flushed
// to disk, then renamed into place. Concurrent writers never share the
// temp file, so a partial write cannot corrupt the real config.
func Save(c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(logging.Dir(), "config-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp) // no-op after a successful rename
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path())
}

// Update performs a locked read-modify-write: it takes the cross-process
// lock, loads the current config, applies fn, and saves, so the GUI and
// the maintenance pass can never lose each other's changes. fn should
// only touch the fields it owns.
func Update(fn func(*Config) error) error {
	unlock, err := acquireLock()
	if err != nil {
		// If locking is unavailable, still perform the update rather than
		// blocking the user; the window is small.
		logging.Logf("config: lock unavailable, proceeding unlocked: %v", err)
		c := Load()
		if err := fn(&c); err != nil {
			return err
		}
		return Save(c)
	}
	defer unlock()
	c := Load()
	if err := fn(&c); err != nil {
		return err
	}
	return Save(c)
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

// RemoveAlert drops a package's alert once it has been dealt with.
func (c *Config) RemoveAlert(pkg string) {
	out := c.Alerts[:0]
	for _, a := range c.Alerts {
		if a.Package != pkg {
			out = append(out, a)
		}
	}
	c.Alerts = out
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

// SetSnapshot records (or clears) the captured pre-apply registry state
// for a fix.
func (c *Config) SetSnapshot(id string, snaps []reg.Snapshot) {
	if c.Snapshots == nil {
		c.Snapshots = map[string][]reg.Snapshot{}
	}
	if len(snaps) == 0 {
		delete(c.Snapshots, id)
		return
	}
	c.Snapshots[id] = snaps
}
