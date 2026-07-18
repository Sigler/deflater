package catalog

import (
	"strings"
	"testing"
)

// The forbidden list backs Deflater's core promise: nothing here may ever
// touch security posture, the gaming stack, driver delivery, or the
// update servicing plumbing. Matched case-insensitively.
var forbiddenPathFragments = []string{
	`\Defender`, `\WindowsDefender`, `SmartScreen`, `SecureBoot`, `\Tpm`, `DeviceGuard`,
	`HypervisorEnforcedCodeIntegrity`, `\Xbox`, `GameBar`, `GameDVR`, `GamingServices`,
	`\WindowsUpdate`, `WindowsSelfHost`, `MemoryManagement`,
	`DriverSearching`, `CurrentControlSet\Services`,
}

var forbiddenPackageFragments = []string{
	"Xbox", "Gaming", "SecHealth", "Defender", "WindowsStore", "DesktopAppInstaller",
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
		// OneDrive fixes may carry no reg/appx: their mechanism is the
		// built-in uninstaller (onedrive-uninstall) or a policy (onedrive-
		// block, which does list a reg op).
		if len(f.Reg) == 0 && len(f.Appx) == 0 && f.Kind != OneDrive {
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

// Some switch values do more than flip a preference: they destroy user
// data (deleting Recall snapshots) and cannot be undone. Those must
// never ride in Light Touch, whose promise is "nothing you would notice
// missing", and must always carry a caution flag.
var dataDestroyingValues = map[string]bool{
	"AllowRecallEnablement": true,
}

func TestDataDestroyingFixesAreGuarded(t *testing.T) {
	for _, f := range Fixes() {
		destroys := false
		for _, op := range f.Reg {
			if dataDestroyingValues[op.Name] {
				destroys = true
			}
		}
		if !destroys {
			continue
		}
		if !f.Caution {
			t.Errorf("%s: destroys data but is not caution-flagged", f.ID)
		}
		for _, p := range f.Profiles {
			if p == LightTouch {
				t.Errorf("%s: destroys data but is in Light Touch", f.ID)
			}
		}
	}
}

func TestRefreshClassification(t *testing.T) {
	// App removals are immediate; policy switches default to a sign-out.
	for _, f := range Fixes() {
		switch f.Kind {
		case AppJunk, AppMight, OneDrive:
			if _, override := refreshOverride[f.ID]; !override && f.RefreshNeeded() != RefreshNone {
				t.Errorf("%s: app removal should refresh immediately, got %q", f.ID, f.RefreshNeeded())
			}
		}
	}
	// Every override id must be a real fix, so the map can't rot.
	for id := range refreshOverride {
		if _, ok := ByID(id); !ok {
			t.Errorf("refreshOverride names unknown fix %q", id)
		}
	}
	// HeaviestRefresh picks the strongest need across a set.
	if got := HeaviestRefresh([]string{"websearch-off", "recall-purge"}); got != RefreshReboot {
		t.Errorf("heaviest of {explorer, reboot} = %q, want reboot", got)
	}
	if got := HeaviestRefresh([]string{"app-news", "websearch-off"}); got != RefreshExplorer {
		t.Errorf("heaviest of {none, explorer} = %q, want explorer", got)
	}
	if got := HeaviestRefresh(nil); got != RefreshNone {
		t.Errorf("heaviest of {} = %q, want none", got)
	}
}
