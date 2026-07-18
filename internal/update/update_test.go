package update

import "testing"

func TestNewerComparesNumerically(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"0.1.2", "0.1.1", true},
		{"0.1.10", "0.1.9", true}, // numeric, not lexical
		{"0.2.0", "0.1.9", true},
		{"1.0.0", "0.9.9", true},
		{"0.1.1", "0.1.1", false}, // equal is not newer
		{"0.1.0", "0.1.1", false}, // older
		{"0.1", "0.1.0", false},   // ragged, equal
		{"0.1.2", "0.1", true},    // ragged, newer patch
		{"", "0.1.0", false},      // empty is not newer
		{"0.1.2-alpha", "0.1.1", true},
		{"0.1.1-alpha", "0.1.1", false}, // suffix ignored, equal
		{"garbage", "0.1.1", false},     // unparseable is not newer
	}
	for _, c := range cases {
		if got := newer(c.a, c.b); got != c.want {
			t.Errorf("newer(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestIsGitHubHost(t *testing.T) {
	good := []string{"github.com", "www.github.com", "api.github.com"}
	bad := []string{"", "evil.com", "github.com.evil.com", "notgithub.com", "githubxcom"}
	for _, h := range good {
		if !isGitHubHost(h) {
			t.Errorf("isGitHubHost(%q) = false, want true", h)
		}
	}
	for _, h := range bad {
		if isGitHubHost(h) {
			t.Errorf("isGitHubHost(%q) = true, want false", h)
		}
	}
}
