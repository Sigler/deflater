package appx

import (
	"reflect"
	"testing"
)

// PowerShell's ConvertTo-Json output changes shape with the result
// count; parseNames must handle all three.
func TestParseNames(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"many packages", `["Microsoft.BingNews","MSTeams"]`, []string{"Microsoft.BingNews", "MSTeams"}},
		{"single package is a bare string", `"Microsoft.BingNews"`, []string{"Microsoft.BingNews"}},
		{"no packages", "", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseNames(c.in)
			if err != nil {
				t.Fatalf("parseNames(%q): %v", c.in, err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("parseNames(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestParseNamesRejectsGarbage(t *testing.T) {
	if _, err := parseNames("not json at all"); err == nil {
		t.Fatal("garbage output must error, not silently return nothing")
	}
}

func TestPrimeControlsInstalled(t *testing.T) {
	s := &Service{}
	s.Prime([]string{"A", "B"})
	installed, err := s.Installed()
	if err != nil {
		t.Fatalf("Installed: %v", err)
	}
	if !installed["A"] || !installed["B"] || installed["C"] {
		t.Fatalf("unexpected installed set: %v", installed)
	}
}
