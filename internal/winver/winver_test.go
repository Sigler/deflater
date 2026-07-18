package winver

import "testing"

func TestFriendlyEdition(t *testing.T) {
	cases := map[string]string{
		"Core":               "Home",
		"CoreSingleLanguage": "Home",
		"Professional":       "Pro",
		"Enterprise":         "Enterprise",
		"Education":          "Education",
		"WeirdFutureSku":     "", // unknown -> no suffix
		"":                   "",
	}
	for id, want := range cases {
		if got := friendlyEdition(id); got != want {
			t.Errorf("friendlyEdition(%q) = %q, want %q", id, got, want)
		}
	}
}

func TestHomeEditionsAreConsumer(t *testing.T) {
	if !homeEditions["Core"] {
		t.Error("Core should be a Home edition")
	}
	if homeEditions["Professional"] || homeEditions["Enterprise"] {
		t.Error("Pro/Enterprise must not be flagged Home")
	}
}

func TestWindowsName(t *testing.T) {
	cases := map[string]string{
		"22631": "Windows 11",
		"26100": "Windows 11",
		"22000": "Windows 11",
		"19045": "Windows 10",
		"0":     "Windows",
		"":      "Windows",
		"junk":  "Windows",
	}
	for build, want := range cases {
		if got := windowsName(build); got != want {
			t.Errorf("windowsName(%q) = %q, want %q", build, got, want)
		}
	}
}
