// Package schtask manages the "Deflater Maintenance" scheduled task.
// The task re-runs Deflater headless (--maintenance) shortly after
// sign-in and weekly, so fixes stay applied after Windows updates and
// the silent-install watcher gets a regular look at the machine.
//
// It runs as the signed-in user with highest privileges, never as
// SYSTEM: many fixes live in the user's own registry hive and app list,
// which SYSTEM would get wrong.
package schtask

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf16"

	"deflater/internal/psrun"
)

const TaskName = "Deflater Maintenance"

const taskXML = `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>Re-applies your Deflater choices after Windows updates and watches for silently installed apps. Never touches Defender, Secure Boot, or anything Xbox.</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <UserId>{USER}</UserId>
      <Delay>PT2M</Delay>
    </LogonTrigger>
    <CalendarTrigger>
      <StartBoundary>2026-01-04T03:00:00</StartBoundary>
      <Enabled>true</Enabled>
      <ScheduleByWeek>
        <DaysOfWeek><Sunday /></DaysOfWeek>
        <WeeksInterval>1</WeeksInterval>
      </ScheduleByWeek>
    </CalendarTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>{USER}</UserId>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <StartWhenAvailable>true</StartWhenAvailable>
    <ExecutionTimeLimit>PT15M</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>{EXE}</Command>
      <Arguments>--maintenance</Arguments>
    </Exec>
  </Actions>
</Task>
`

func currentUser() string {
	return os.Getenv("USERDOMAIN") + `\` + os.Getenv("USERNAME")
}

func xmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return r.Replace(s)
}

// buildTaskXML fills the task template for the given user and target.
func buildTaskXML(user, exePath string) string {
	xml := strings.ReplaceAll(taskXML, "{USER}", xmlEscape(user))
	return strings.ReplaceAll(xml, "{EXE}", xmlEscape(exePath))
}

// Install registers (or replaces) the task pointing at exePath.
func Install(exePath string) error {
	xml := buildTaskXML(currentUser(), exePath)

	// Task Scheduler expects task XML files in UTF-16 LE. A uniquely-named
	// temp file avoids a predictable path a local attacker could swap
	// between our write and schtasks reading it.
	f, err := os.CreateTemp("", "deflater-task-*.xml")
	if err != nil {
		return err
	}
	tmp := f.Name()
	f.Close()
	defer os.Remove(tmp)
	if err := os.WriteFile(tmp, utf16LE(xml), 0o600); err != nil {
		return err
	}

	_, err = psrun.Run(fmt.Sprintf(
		`schtasks.exe /Create /F /TN "%s" /XML "%s"`, TaskName, tmp), 60*time.Second)
	return err
}

// Uninstall removes the task. Removing a task that does not exist is fine;
// this is decided by querying first, so it works on any UI language
// rather than matching an English error string.
func Uninstall() error {
	if !IsInstalled() {
		return nil
	}
	_, err := psrun.Run(fmt.Sprintf(`schtasks.exe /Delete /F /TN "%s"`, TaskName), 60*time.Second)
	return err
}

// IsInstalled reports whether the task exists.
func IsInstalled() bool {
	_, err := psrun.Run(fmt.Sprintf(`schtasks.exe /Query /TN "%s" | Out-Null`, TaskName), 60*time.Second)
	return err == nil
}

func utf16LE(s string) []byte {
	codes := utf16.Encode([]rune(s))
	buf := make([]byte, 0, 2+len(codes)*2)
	buf = append(buf, 0xFF, 0xFE) // BOM
	for _, c := range codes {
		buf = append(buf, byte(c), byte(c>>8))
	}
	return buf
}
