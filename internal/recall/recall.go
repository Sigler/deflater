// Package recall reports whether Windows Recall has stored snapshots on
// this PC, and how much space they take. It only reads: it never opens,
// decodes, or deletes the snapshots. Removing them is the job of the
// recall-purge fix; this is the privacy signal shown before the user
// decides.
package recall

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Info describes the Recall snapshot store for the current user.
type Info struct {
	// Present is true only when the store exists and holds data.
	Present bool `json:"present"`
	// Path is the snapshot store's folder (shown so the user can find it).
	Path string `json:"path"`
	// Bytes is the total size on disk of everything under Path.
	Bytes int64 `json:"bytes"`
}

// storeRoot is Recall's per-user snapshot store. The UKP folder under
// CoreAIPlatform.00 holds the screenshot ImageStore and its database.
func storeRoot() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "CoreAIPlatform.00", "UKP")
}

// Detect walks the store and returns its size. Best-effort: unreadable
// entries are skipped rather than failing the whole scan, and a missing
// store simply reports Present=false.
func Detect() Info {
	root := storeRoot()
	fi, err := os.Stat(root)
	if err != nil || !fi.IsDir() {
		return Info{Path: root}
	}
	var total int64
	_ = filepath.WalkDir(root, func(_ string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}
		if info, e := d.Info(); e == nil {
			total += info.Size()
		}
		return nil
	})
	return Info{Present: total > 0, Path: root, Bytes: total}
}

// IsStorePath reports whether p is exactly the Recall store folder, so a
// caller (e.g. an "open folder" action) can refuse any other path.
func IsStorePath(p string) bool {
	return p == storeRoot()
}
