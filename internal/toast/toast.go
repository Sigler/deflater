// Package toast shows native Windows notifications. Deflater is an
// unpackaged app, so it registers an AppUserModelID once under HKCU and
// raises toasts through the WinRT API via PowerShell.
package toast

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	"deflater/internal/psrun"
)

const aumid = "Deflater.App"

// Register declares the app identity toasts are attributed to. Safe to
// call every launch; it writes two small HKCU values.
func Register(iconPath string) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Classes\AppUserModelId\`+aumid, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if err := k.SetStringValue("DisplayName", "Deflater"); err != nil {
		return err
	}
	if iconPath != "" {
		if err := k.SetStringValue("IconUri", iconPath); err != nil {
			return err
		}
	}
	return nil
}

func psQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// Show raises a toast with a title line and a body line.
func Show(title, body string) error {
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType=WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType=WindowsRuntime] | Out-Null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$toast = @'
<toast><visual><binding template="ToastGeneric"><text>{TITLE}</text><text>{BODY}</text></binding></visual></toast>
'@
$toast = $toast.Replace('{TITLE}', %s).Replace('{BODY}', %s)
$xml.LoadXml($toast)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier(%s).Show(
    (New-Object Windows.UI.Notifications.ToastNotification $xml))
`, psQuote(escapeXML(title)), psQuote(escapeXML(body)), psQuote(aumid))
	_, err := psrun.Run(script, 30*time.Second)
	return err
}

func escapeXML(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return r.Replace(s)
}
