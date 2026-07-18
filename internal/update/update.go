// Package update provides best-effort, fail-silent update awareness: it
// asks GitHub for the latest published release and reports whether it is
// newer than this build. It never downloads or installs anything; the UI
// just shows a link to the releases page when something newer exists.
package update

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const releasesAPI = "https://api.github.com/repos/Sigler/deflater/releases/latest"

// ReleasesPage is the human-facing releases list, used as the link target
// and as a fallback when the API doesn't hand back a specific URL.
const ReleasesPage = "https://github.com/Sigler/deflater/releases"

// Info is the result of a check. Available is true only when a strictly
// newer release exists; the UI keys off it.
type Info struct {
	Available bool   `json:"available"`
	Current   string `json:"current"`
	Latest    string `json:"latest"`
	URL       string `json:"url"`
}

// Check asks GitHub for the latest release and compares it to current.
// Fail-silent by design: any network, HTTP, or parse error returns
// Info{Available:false} with the releases page as the URL, so update
// awareness can never disrupt the app or block startup.
func Check(current string) Info {
	// Strip a leading "v" symmetrically with the tag name below, so a
	// future appVersion of "v0.2.0" can't read as major 0 and nag.
	current = strings.TrimPrefix(strings.TrimSpace(current), "v")
	info := Info{Current: current, URL: ReleasesPage}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, releasesAPI, nil)
	if err != nil {
		return info
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "Deflater")

	resp, err := client.Do(req)
	if err != nil {
		return info
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return info
	}

	var payload struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}
	// Cap the body: a hostile or broken endpoint shouldn't be able to make
	// us allocate unboundedly. A release JSON is a few KB.
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return info
	}

	info.Latest = strings.TrimPrefix(strings.TrimSpace(payload.TagName), "v")
	// The link is opened via the OS shell, so only trust an https github.com
	// URL from the response; otherwise fall back to the hardcoded page.
	if u, err := url.Parse(payload.HTMLURL); err == nil && u.Scheme == "https" && isGitHubHost(u.Host) {
		info.URL = payload.HTMLURL
	}
	info.Available = newer(info.Latest, current)
	return info
}

// isGitHubHost reports whether host is github.com or a subdomain of it.
func isGitHubHost(host string) bool {
	return host == "github.com" || strings.HasSuffix(host, ".github.com")
}

// newer reports whether version a is strictly greater than b, comparing
// dot-separated numeric components (so 0.1.10 > 0.1.9). Missing or
// non-numeric components count as 0, and a tie or any parse trouble means
// "not newer", so we never nag on a bad or equal version string.
func newer(a, b string) bool {
	as := strings.Split(a, ".")
	bs := strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		ai, bi := component(as, i), component(bs, i)
		if ai != bi {
			return ai > bi
		}
	}
	return false
}

// component returns the numeric value of the i-th dot-part, tolerating a
// pre-release suffix (1.2.0-alpha -> 2 at index 2) and treating anything
// missing or unparseable as 0.
func component(parts []string, i int) int {
	if i >= len(parts) {
		return 0
	}
	p := parts[i]
	if cut := strings.IndexAny(p, "-+"); cut >= 0 {
		p = p[:cut]
	}
	n, err := strconv.Atoi(strings.TrimSpace(p))
	if err != nil {
		return 0
	}
	return n
}
