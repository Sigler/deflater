// Package logging writes a plain append-only log file under
// %LOCALAPPDATA%\Deflater\logs so users (and we) can see what
// Deflater actually did and when.
package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const maxLogBytes = 1 << 20 // rotate after 1 MB, keeps one previous file

var mu sync.Mutex

// Dir returns the Deflater data directory, creating it if needed.
// DEFLATER_DATA_DIR overrides the location (tests, portable use).
func Dir() string {
	d := os.Getenv("DEFLATER_DATA_DIR")
	if d == "" {
		d = filepath.Join(os.Getenv("LOCALAPPDATA"), "Deflater")
	}
	_ = os.MkdirAll(d, 0o755)
	return d
}

// LogDir returns the log directory, creating it if needed.
func LogDir() string {
	d := filepath.Join(Dir(), "logs")
	_ = os.MkdirAll(d, 0o755)
	return d
}

func logPath() string { return filepath.Join(LogDir(), "deflater.log") }

// Logf appends one timestamped line to the log file.
func Logf(format string, args ...any) {
	mu.Lock()
	defer mu.Unlock()

	p := logPath()
	if info, err := os.Stat(p); err == nil && info.Size() > maxLogBytes {
		_ = os.Rename(p, p+".old")
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	line := fmt.Sprintf("%s  %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))
	_, _ = f.WriteString(line)
}
