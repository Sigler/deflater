// Package catalog is the single source of truth for what Deflater can
// fix: the exact registry values and Store packages behind every switch,
// which category each fix lives in, and which profiles preselect it.
//
// Ground rules for anything added here:
//   - Reversible: registry switches undo with one click (deleting the
//     value or restoring the Settings default), app removals reinstall
//     from the Microsoft Store. Apps that are no longer reinstallable
//     (Maps, People, Movies & TV) are deliberately absent.
//   - Never touched: Defender, Secure Boot, TPM, VBS, driver delivery,
//     Windows Update servicing, or anything Xbox / Game Pass. Anti-cheat
//     sees a stock security posture, always.
//
// Every mechanism here was verified against Microsoft's policy
// documentation for Windows 11 24H2/25H2 in July 2026. Notable
// verification results are commented on the fixes they concern.
//
// User-facing names and descriptions live in the frontend string catalog
// (frontend/src/lib/strings), keyed by these ids, so more languages can
// be added without touching Go.
package catalog

import (
	"fmt"

	"deflater/internal/reg"
)

// Kind describes how a fix operates.
type Kind string

const (
	Switch   Kind = "switch"    // registry values, fully reversible in-app
	AppJunk  Kind = "app-junk"  // Store app nobody installs on purpose
	AppMight Kind = "app-might" // Store app some people genuinely use
	OneDrive Kind = "onedrive"  // policy switch plus Microsoft's own uninstaller
)

// Profile ids, mild to maximum.
const (
	LightTouch  = "light-touch"
	CleanSweep  = "clean-sweep"
	FullDeflate = "full-deflate"
)

// Categories lists category ids in display order.
var Categories = []string{
	"ads-nags",
	"junk-apps",
	"start-search",
	"copilot-ai",
	"privacy",
	"might-use",
}

// Fix is one toggleable improvement.
type Fix struct {
	ID       string `json:"id"`
	Category string `json:"category"`
	Kind     Kind   `json:"kind"`
	// Caution marks fixes that remove or hide something a person might
	// use on purpose; the UI renders these for careful review.
	Caution  bool     `json:"caution"`
	Profiles []string `json:"profiles"`
	Reg      []reg.Op `json:"reg,omitempty"`
	Appx     []string `json:"appx,omitempty"`
	// Group ties related fixes together in the UI so a bundled action can
	// be split into independent toggles (e.g. block OneDrive vs also
	// uninstall it) that share one card. Empty means a standalone fix.
	Group string `json:"group,omitempty"`
	// Primary marks the lead fix of a Group; the UI renders it first and
	// treats the rest as secondary sub-options.
	Primary bool `json:"primary,omitempty"`
}

// Refresh is the lightest action that makes a fix visibly take effect.
type Refresh string

const (
	RefreshNone     Refresh = "none"     // live the moment it's applied
	RefreshExplorer Refresh = "explorer" // restart Explorer (taskbar/Start/Search/File Explorer)
	RefreshSignOut  Refresh = "signout"  // sign out and back in
	RefreshReboot   Refresh = "reboot"   // full restart
)

// refreshOverride names the fixes whose refresh need differs from the
// per-Kind default. Kept as a map so the catalog entries stay terse.
//   - The shell-surface switches only need a shell restart, not a full
//     sign-out. RestartExplorer restarts explorer.exe AND SearchHost.exe
//     (the search box, for websearch-off / search-highlights) and
//     StartMenuExperienceHost.exe (Start recommendations, for
//     settings-suggestions), which is what actually re-reads these.
//   - recall-purge removes a component and needs a reboot to complete.
var refreshOverride = map[string]Refresh{
	"websearch-off":        RefreshExplorer,
	"widgets":              RefreshExplorer,
	"search-highlights":    RefreshExplorer,
	"explorer-ads":         RefreshExplorer,
	"settings-suggestions": RefreshExplorer,
	"recall-purge":         RefreshReboot,
}

// RefreshNeeded reports the lightest action that makes this fix take
// effect: nothing for app removals (gone immediately), the override where
// one is listed, and a sign-out as the safe default for policy switches.
func (f Fix) RefreshNeeded() Refresh {
	if r, ok := refreshOverride[f.ID]; ok {
		return r
	}
	switch f.Kind {
	case AppJunk, AppMight, OneDrive:
		return RefreshNone
	}
	return RefreshSignOut
}

func refreshRank(r Refresh) int {
	switch r {
	case RefreshNone:
		return 0
	case RefreshExplorer:
		return 1
	case RefreshReboot:
		return 3
	default: // RefreshSignOut and anything unknown
		return 2
	}
}

// HeaviestRefresh returns the strongest refresh any of the given fix ids
// needs, so the UI can name the single action that covers them all.
func HeaviestRefresh(ids []string) Refresh {
	heaviest := RefreshNone
	for _, id := range ids {
		f, ok := ByID(id)
		if !ok {
			continue
		}
		if r := f.RefreshNeeded(); refreshRank(r) > refreshRank(heaviest) {
			heaviest = r
		}
	}
	return heaviest
}

func all() []string   { return []string{LightTouch, CleanSweep, FullDeflate} }
func sweep() []string { return []string{CleanSweep, FullDeflate} }
func full() []string  { return []string{FullDeflate} }

// pol is a policy value: absent by default, so revert deletes it.
func pol(hive, path, name string, value uint32) reg.Op {
	return reg.Op{Hive: hive, Path: path, Name: name, Value: value, Revert: "delete"}
}

// tog is a Settings-backed toggle with a known default: revert restores
// the default explicitly so the Settings UI reads correctly again. (When
// a value was captured at apply time the engine restores that instead,
// which is more accurate; this static default is the fallback.)
func tog(hive, path, name string, value, defaultValue uint32) reg.Op {
	return reg.Op{Hive: hive, Path: path, Name: name, Value: value, Revert: fmt.Sprintf("set:%d", defaultValue)}
}

const (
	cloudContent = `SOFTWARE\Policies\Microsoft\Windows\CloudContent`
	cdm          = `Software\Microsoft\Windows\CurrentVersion\ContentDeliveryManager`
	explorerAdv  = `Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced`
	edgePol      = `SOFTWARE\Policies\Microsoft\Edge`
	windowsAI    = `SOFTWARE\Policies\Microsoft\Windows\WindowsAI`
)

// Fixes returns the full catalog in display order within each category.
func Fixes() []Fix {
	return []Fix{
		// ---- Ads and nags -------------------------------------------------
		{
			ID: "lockscreen-ads", Category: "ads-nags", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				tog("HKCU", cdm, "RotatingLockScreenOverlayEnabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-338387Enabled", 0, 1),
			},
		},
		{
			ID: "explorer-ads", Category: "ads-nags", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				tog("HKCU", explorerAdv, "ShowSyncProviderNotifications", 0, 1),
			},
		},
		{
			ID: "scoobe-off", Category: "ads-nags", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				tog("HKCU", `Software\Microsoft\Windows\CurrentVersion\UserProfileEngagement`, "ScoobeSystemSettingEnabled", 0, 1),
			},
		},
		{
			ID: "suggested-toasts-off", Category: "ads-nags", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKCU", `Software\Microsoft\Windows\CurrentVersion\Notifications\Settings\Windows.SystemToast.Suggested`, "Enabled", 0),
			},
		},
		{
			ID: "settings-suggestions", Category: "ads-nags", Kind: Switch, Profiles: all(),
			// 338393/353694/353696 together are "suggested content in
			// Settings"; 338389 is tips, 310093 the post-update welcome.
			Reg: []reg.Op{
				tog("HKCU", cdm, "SystemPaneSuggestionsEnabled", 0, 1),
				tog("HKCU", cdm, "SoftLandingEnabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-338388Enabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-338389Enabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-338393Enabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-353694Enabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-353696Enabled", 0, 1),
				tog("HKCU", cdm, "SubscribedContent-310093Enabled", 0, 1),
				tog("HKCU", explorerAdv, "Start_IrisRecommendations", 0, 1),
			},
		},
		{
			// Note: any Edge policy makes Edge show "Managed by your
			// organization" in its menu; the UI copy discloses this.
			ID: "edge-nags", Category: "ads-nags", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", edgePol, "StartupBoostEnabled", 0),
				pol("HKLM", edgePol, "BackgroundModeEnabled", 0),
				pol("HKLM", edgePol, "HubsSidebarEnabled", 0),
				pol("HKLM", edgePol, "ShowRecommendationsEnabled", 0),
				pol("HKLM", edgePol, "PersonalizationReportingEnabled", 0),
			},
		},
		{
			// RemoveDesktopShortcutDefault=1 also deletes an existing Edge
			// desktop shortcut, which is why this is not in Light Touch.
			ID: "edge-shortcut", Category: "ads-nags", Kind: Switch, Caution: true, Profiles: sweep(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\EdgeUpdate`, "CreateDesktopShortcutDefault", 0),
				pol("HKLM", `SOFTWARE\Policies\Microsoft\EdgeUpdate`, "RemoveDesktopShortcutDefault", 1),
			},
		},

		// ---- Junk apps (and stopping new ones) ----------------------------
		{
			// The HKCU ContentDeliveryManager values are the effective
			// switch on Pro/Home; the HKLM CloudContent policies are
			// documented Enterprise/Education but kept as harmless,
			// future-proof hardening.
			ID: "silent-app-installs", Category: "junk-apps", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				tog("HKCU", cdm, "SilentInstalledAppsEnabled", 0, 1),
				tog("HKCU", cdm, "PreInstalledAppsEnabled", 0, 1),
				tog("HKCU", cdm, "OemPreInstalledAppsEnabled", 0, 1),
				pol("HKLM", cloudContent, "DisableWindowsConsumerFeatures", 1),
				pol("HKLM", cloudContent, "DisableConsumerAccountStateContent", 1),
				pol("HKLM", cloudContent, "DisableCloudOptimizedContent", 1),
				pol("HKLM", cloudContent, "DisableSoftLanding", 1),
			},
		},
		{
			// The policy key overrides the Settings UI; the CurrentVersion
			// value is what the "Device installation settings" switch
			// writes. Both are set so the UI reflects reality. Driver
			// delivery (DriverSearching) is deliberately not touched.
			ID: "device-metadata-off", Category: "junk-apps", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\Device Metadata`, "PreventDeviceMetadataFromNetwork", 1),
				tog("HKLM", `SOFTWARE\Microsoft\Windows\CurrentVersion\Device Metadata`, "PreventDeviceMetadataFromNetwork", 1, 0),
			},
		},
		{ID: "app-officehub", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.MicrosoftOfficeHub"}},
		{ID: "app-news", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.BingNews"}},
		{ID: "app-weather", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.BingWeather"}},
		{ID: "app-solitaire", Category: "junk-apps", Kind: AppJunk, Caution: true, Profiles: sweep(), Appx: []string{"Microsoft.MicrosoftSolitaireCollection"}},
		{ID: "app-gethelp", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.GetHelp"}},
		{ID: "app-feedback", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.WindowsFeedbackHub"}},
		{ID: "app-bingsearch", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.BingSearch"}},
		{ID: "app-powerautomate", Category: "junk-apps", Kind: AppJunk, Profiles: sweep(), Appx: []string{"Microsoft.PowerAutomateDesktop"}},

		// ---- Start menu, search and taskbar -------------------------------
		{
			ID: "websearch-off", Category: "start-search", Kind: Switch, Caution: true, Profiles: sweep(),
			Reg: []reg.Op{
				pol("HKCU", `Software\Policies\Microsoft\Windows\Explorer`, "DisableSearchBoxSuggestions", 1),
			},
		},
		{
			// Policy disables the feature machine-wide; TaskbarDa hides
			// the button for this user so the taskbar updates cleanly.
			ID: "widgets", Category: "start-search", Kind: Switch, Caution: true, Profiles: sweep(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Dsh`, "AllowNewsAndInterests", 0),
				tog("HKCU", explorerAdv, "TaskbarDa", 0, 1),
			},
		},
		{
			ID: "search-highlights", Category: "start-search", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\Windows Search`, "EnableDynamicContentInWSB", 0),
				tog("HKCU", `Software\Microsoft\Windows\CurrentVersion\SearchSettings`, "IsDynamicSearchBoxEnabled", 0, 1),
			},
		},

		// ---- Copilot and AI ----------------------------------------------
		{
			// The old TurnOffWindowsCopilot policy is deprecated and does
			// not govern the Copilot store app on 24H2+; removing the app
			// is the mechanism that actually works, so it lives in Clean
			// Sweep where "Copilot off" belongs.
			ID: "app-copilot", Category: "copilot-ai", Kind: AppMight, Caution: true, Profiles: sweep(), Appx: []string{"Microsoft.Copilot"},
		},
		{
			// Stops Recall from saving new snapshots. Harmless and fully
			// reversible, so it stays in every profile.
			ID: "recall-off", Category: "copilot-ai", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", windowsAI, "DisableAIDataAnalysis", 1),
				pol("HKCU", `Software\Policies\Microsoft\Windows\WindowsAI`, "DisableAIDataAnalysis", 1),
			},
		},
		{
			// AllowRecallEnablement=0 removes the Recall component AND
			// permanently deletes any existing snapshots. That data loss
			// is why this is a separate, caution-flagged fix kept out of
			// Light Touch, unlike the reversible snapshot pause above.
			ID: "recall-purge", Category: "copilot-ai", Kind: Switch, Caution: true, Profiles: sweep(),
			Reg: []reg.Op{
				pol("HKLM", windowsAI, "AllowRecallEnablement", 0),
			},
		},
		{
			ID: "click-to-do-off", Category: "copilot-ai", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", windowsAI, "DisableClickToDo", 1),
			},
		},

		// ---- Privacy ------------------------------------------------------
		{
			ID: "advertising-id", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\AdvertisingInfo`, "DisabledByGroupPolicy", 1),
				tog("HKCU", `Software\Microsoft\Windows\CurrentVersion\AdvertisingInfo`, "Enabled", 0, 1),
			},
		},
		{
			// On Pro the minimum is Required (1); 0 is silently treated as
			// 1, so the UI says "minimum", never "off".
			ID: "telemetry-minimum", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\DataCollection`, "AllowTelemetry", 1),
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\DataCollection`, "DoNotShowFeedbackNotifications", 1),
			},
		},
		{
			// This policy is user-scope: it must live in HKCU. The HKLM
			// copy many scripts write is a silent no-op.
			ID: "tailored-experiences", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKCU", `Software\Policies\Microsoft\Windows\CloudContent`, "DisableTailoredExperiencesWithDiagnosticData", 1),
				tog("HKCU", `Software\Microsoft\Windows\CurrentVersion\Privacy`, "TailoredExperiencesWithDiagnosticDataEnabled", 0, 1),
			},
		},
		{
			ID: "activity-history", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\System`, "PublishUserActivities", 0),
			},
		},
		{
			ID: "inking-personalization", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				tog("HKCU", `Software\Microsoft\InputPersonalization`, "RestrictImplicitInkCollection", 1, 0),
				tog("HKCU", `Software\Microsoft\InputPersonalization`, "RestrictImplicitTextCollection", 1, 0),
				tog("HKCU", `Software\Microsoft\InputPersonalization\TrainedDataStore`, "HarvestContacts", 0, 1),
				tog("HKCU", `Software\Microsoft\Personalization\Settings`, "AcceptedPrivacyPolicy", 0, 1),
			},
		},
		{
			// Mode 1 keeps LAN peering and normal updates; it only stops
			// internet peer upload/download.
			ID: "delivery-optimization", Category: "privacy", Kind: Switch, Profiles: all(),
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\DeliveryOptimization`, "DODownloadMode", 1),
			},
		},

		// ---- Apps you might use ------------------------------------------
		{
			// The reversible half: a policy that stops OneDrive running and
			// nagging, while leaving the app installed. The primary of the
			// OneDrive group.
			ID: "onedrive-block", Category: "might-use", Kind: Switch, Caution: true, Profiles: full(),
			Group: "onedrive", Primary: true,
			Reg: []reg.Op{
				pol("HKLM", `SOFTWARE\Policies\Microsoft\Windows\OneDrive`, "DisableFileSyncNGSC", 1),
			},
		},
		{
			// The drastic half: run Microsoft's own uninstaller. No registry
			// of its own; a secondary sub-option under onedrive-block. Cloud
			// files are untouched, and Microsoft.OneDriveSync (a sync
			// component, not the app) is deliberately NOT removed.
			ID: "onedrive-uninstall", Category: "might-use", Kind: OneDrive, Caution: true, Profiles: full(),
			Group: "onedrive",
		},
		{
			// Phone Link and the Cross Device Experience Host travel
			// together; removing one without the other leaves phone
			// integration half-broken.
			ID: "app-phonelink", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(),
			Appx: []string{"Microsoft.YourPhone", "MicrosoftWindows.CrossDevice"},
		},
		{ID: "app-teams", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"MSTeams"}},
		{ID: "app-outlook", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"Microsoft.OutlookForWindows"}},
		{ID: "app-clipchamp", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"Clipchamp.Clipchamp"}},
		{ID: "app-todo", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"Microsoft.Todos"}},
		{ID: "app-family", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"MicrosoftCorporationII.MicrosoftFamily"}},
		{ID: "app-quickassist", Category: "might-use", Kind: AppMight, Caution: true, Profiles: full(), Appx: []string{"MicrosoftCorporationII.QuickAssist"}},
	}
}

// ByID returns the fix with the given id, or false.
func ByID(id string) (Fix, bool) {
	for _, f := range Fixes() {
		if f.ID == id {
			return f, true
		}
	}
	return Fix{}, false
}

// ManagedPackages returns every package name managed by the given fix
// ids; the watcher uses this to avoid alerting on packages Deflater
// itself removes or re-removes.
func ManagedPackages(ids []string) map[string]bool {
	out := map[string]bool{}
	for _, id := range ids {
		if f, ok := ByID(id); ok {
			for _, p := range f.Appx {
				out[p] = true
			}
		}
	}
	return out
}
