package schtask

import (
	"strings"
	"testing"
)

func TestTaskXMLStructure(t *testing.T) {
	xml := buildTaskXML(`DESKTOP\mike`, `C:\Users\mike\AppData\Local\Deflater\bin\Deflater.exe`)

	checks := []string{
		"<Delay>PT2M</Delay>",
		"<Sunday />",
		"<RunLevel>HighestAvailable</RunLevel>",
		"<LogonType>InteractiveToken</LogonType>",
		"<Arguments>--maintenance</Arguments>",
		`DESKTOP\mike`,
		`C:\Users\mike\AppData\Local\Deflater\bin\Deflater.exe`,
		"<ExecutionTimeLimit>PT15M</ExecutionTimeLimit>",
	}
	for _, want := range checks {
		if !strings.Contains(xml, want) {
			t.Errorf("task XML missing %q", want)
		}
	}
	if strings.Contains(xml, "{USER}") || strings.Contains(xml, "{EXE}") {
		t.Error("template placeholders were not replaced")
	}
}

func TestTaskXMLEscapesSpecialCharacters(t *testing.T) {
	xml := buildTaskXML("DOM&AIN\\bob", `C:\odd & path\Deflater.exe`)
	if !strings.Contains(xml, "DOM&amp;AIN") {
		t.Error("user with & must be XML-escaped")
	}
	if !strings.Contains(xml, `C:\odd &amp; path\Deflater.exe`) {
		t.Error("path with & must be XML-escaped")
	}
	if strings.Contains(xml, "odd & path") {
		t.Error("raw & survived into the XML")
	}
}

func TestForeignTaskRecognition(t *testing.T) {
	// The vetted list must recognize its own entries case-insensitively.
	for _, ft := range knownForeign {
		if !IsKnownForeign(ft.Name) {
			t.Errorf("known task %q not recognized", ft.Name)
		}
		if !IsKnownForeign(strings.ToUpper(ft.Name)) {
			t.Errorf("known task %q not recognized case-insensitively", ft.Name)
		}
	}
	// Anything not on the list must be refused, so a caller can never
	// drive removal of an arbitrary task.
	for _, bad := range []string{"", "OneDrive Standalone Update Task", `\Microsoft\Windows\Foo`, "Deflater Maintenance"} {
		if IsKnownForeign(bad) {
			t.Errorf("unlisted task %q was recognized", bad)
		}
		if err := RemoveForeign(bad); err == nil {
			t.Errorf("RemoveForeign(%q) should refuse an unrecognized task", bad)
		}
	}
}

func TestUTF16LEHasBOMAndLittleEndian(t *testing.T) {
	b := utf16LE("ab")
	want := []byte{0xFF, 0xFE, 'a', 0x00, 'b', 0x00}
	if len(b) != len(want) {
		t.Fatalf("length %d, want %d", len(b), len(want))
	}
	for i := range want {
		if b[i] != want[i] {
			t.Fatalf("byte %d = %x, want %x", i, b[i], want[i])
		}
	}
}
