package catalog

import (
	"strings"
	"testing"
)

// The forbidden list backs Deflater's core promise: nothing here may ever
// touch security posture or the gaming stack.
var forbiddenPathFragments = []string{
	`\Defender`, `\WindowsDefender`, `SecureBoot`, `\Tpm`, `DeviceGuard`,
	`HypervisorEnforcedCodeIntegrity`, `\Xbox`, `GameBar`, `GamingServices`,
	`\WindowsUpdate`, `MemoryManagement`,
}

var forbiddenPackageFragments = []string{
	"Xbox", "Gaming", "SecHealth", "Defender",
}

func TestCatalogIntegrity(t *testing.T) {
	seen := map[string]bool{}
	validCats := map[string]bool{}
	for _, c := range Categories {
		validCats[c] = true
	}
	validProfiles := map[string]bool{LightTouch: true, CleanSweep: true, FullDeflate: true}

	for _, f := range Fixes() {
		if f.ID == "" {
			t.Fatal("fix with empty id")
		}
		if seen[f.ID] {
			t.Fatalf("duplicate fix id %q", f.ID)
		}
		seen[f.ID] = true

		if !validCats[f.Category] {
			t.Errorf("%s: unknown category %q", f.ID, f.Category)
		}
		if len(f.Profiles) == 0 {
			t.Errorf("%s: belongs to no profile", f.ID)
		}
		for _, p := range f.Profiles {
			if !validProfiles[p] {
				t.Errorf("%s: unknown profile %q", f.ID, p)
			}
		}
		if len(f.Reg) == 0 && len(f.Appx) == 0 {
			t.Errorf("%s: has no mechanism", f.ID)
		}
		switch f.Kind {
		case Switch:
			if len(f.Appx) != 0 {
				t.Errorf("%s: switches must not remove apps", f.ID)
			}
		case AppJunk, AppMight:
			if len(f.Reg) != 0 {
				t.Errorf("%s: app removals must not write registry values", f.ID)
			}
			if len(f.Appx) == 0 {
				t.Errorf("%s: app removal without package names", f.ID)
			}
		case OneDrive:
			// combined mechanism, checked by the forbidden scans below
		default:
			t.Errorf("%s: unknown kind %q", f.ID, f.Kind)
		}
	}
}

func TestNothingForbiddenIsTouched(t *testing.T) {
	for _, f := range Fixes() {
		for _, op := range f.Reg {
			if op.Hive != "HKLM" && op.Hive != "HKCU" {
				t.Errorf("%s: bad hive %q", f.ID, op.Hive)
			}
			for _, frag := range forbiddenPathFragments {
				if strings.Contains(strings.ToLower(op.Path), strings.ToLower(frag)) {
					t.Errorf("%s: registry path %q touches forbidden area %q", f.ID, op.Path, frag)
				}
			}
		}
		for _, pkg := range f.Appx {
			for _, frag := range forbiddenPackageFragments {
				if strings.Contains(strings.ToLower(pkg), strings.ToLower(frag)) {
					t.Errorf("%s: package %q touches forbidden area %q", f.ID, pkg, frag)
				}
			}
		}
	}
}

func TestProfilesNest(t *testing.T) {
	// The profiles are an intensity scale: everything in Light Touch must
	// be in Clean Sweep, and everything in Clean Sweep in Full Deflate.
	in := func(f Fix, profile string) bool {
		for _, p := range f.Profiles {
			if p == profile {
				return true
			}
		}
		return false
	}
	for _, f := range Fixes() {
		if in(f, LightTouch) && !in(f, CleanSweep) {
			t.Errorf("%s: in Light Touch but missing from Clean Sweep", f.ID)
		}
		if in(f, CleanSweep) && !in(f, FullDeflate) {
			t.Errorf("%s: in Clean Sweep but missing from Full Deflate", f.ID)
		}
	}
}

func TestLightTouchRemovesNothing(t *testing.T) {
	for _, f := range Fixes() {
		for _, p := range f.Profiles {
			if p == LightTouch && f.Kind != Switch {
				t.Errorf("%s: Light Touch must only flip switches, found kind %q", f.ID, f.Kind)
			}
		}
	}
}
