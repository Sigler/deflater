// Package winver reports the Windows edition, so the UI can name it and
// (later) tailor which fixes it highlights. It reads CurrentVersion, the
// same place Settings and gpedit read from.
package winver

import (
	"strconv"

	"golang.org/x/sys/windows/registry"
)

// Info describes the running Windows edition.
type Info struct {
	// Edition is a friendly label like "Windows 11 Home".
	Edition string `json:"edition"`
	// Home is true for the consumer editions (Core*), where gpedit.msc
	// doesn't exist, which is exactly why Deflater writes the registry
	// directly.
	Home bool `json:"home"`
}

// homeEditions are the consumer EditionID values that lack Group Policy.
var homeEditions = map[string]bool{
	"Core":                true,
	"CoreN":               true,
	"CoreSingleLanguage":  true,
	"CoreCountrySpecific": true,
}

// friendlyEdition turns an EditionID into a display word, or "" if it's
// one we don't have a nicer name for (then only "Windows 11" is shown).
func friendlyEdition(id string) string {
	switch id {
	case "Core", "CoreN", "CoreSingleLanguage", "CoreCountrySpecific":
		return "Home"
	case "Professional", "ProfessionalN":
		return "Pro"
	case "ProfessionalWorkstation", "ProfessionalWorkstationN":
		return "Pro for Workstations"
	case "Enterprise", "EnterpriseN":
		return "Enterprise"
	case "Education", "EducationN":
		return "Education"
	case "ProfessionalEducation", "ProfessionalEducationN":
		return "Pro Education"
	case "ServerStandard", "ServerDatacenter":
		return "Server"
	default:
		return ""
	}
}

// windowsName returns "Windows 11", "Windows 10", or a plain "Windows"
// fallback, derived from the build number (ProductName lags and still
// says "Windows 10" on 11).
func windowsName(currentBuild string) string {
	n, err := strconv.Atoi(currentBuild)
	if err != nil || n == 0 {
		return "Windows"
	}
	if n >= 22000 {
		return "Windows 11"
	}
	return "Windows 10"
}

// Detect reads the edition. On any read failure it degrades to a plain
// "Windows" label and Home=false, never an error.
func Detect() Info {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return Info{Edition: "Windows"}
	}
	defer k.Close()

	editionID, _, _ := k.GetStringValue("EditionID")
	currentBuild, _, _ := k.GetStringValue("CurrentBuild")

	edition := windowsName(currentBuild)
	if word := friendlyEdition(editionID); word != "" {
		edition += " " + word
	}
	return Info{Edition: edition, Home: homeEditions[editionID]}
}
