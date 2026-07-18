package catalog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every fix in the catalog must have user-facing copy in the frontend
// string catalog, or it would render as a bare id.
func TestEveryFixHasStrings(t *testing.T) {
	path := filepath.Join("..", "..", "frontend", "src", "lib", "strings", "en.ts")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	content := string(data)
	for _, f := range Fixes() {
		// Keys appear either quoted ("app-news":) or bare (widgets:).
		if !strings.Contains(content, `"`+f.ID+`"`) && !strings.Contains(content, f.ID+":") {
			t.Errorf("fix %q has no entry in frontend strings/en.ts", f.ID)
		}
	}
	for _, c := range Categories {
		if !strings.Contains(content, `"`+c+`"`) && !strings.Contains(content, c+":") {
			t.Errorf("category %q has no entry in frontend strings/en.ts", c)
		}
	}
}
