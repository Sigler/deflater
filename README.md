<p align="center">
  <img src="assets/mascot.png" width="440" alt="Four worn blue vinyl cushions arranged like the Windows logo, one deflated, wrapped in a red DEFLATED: ANTI-BLOAT ENFORCEMENT banner" />
</p>

<h1 align="center">Deflater</h1>
<p align="center">Let the air out of Windows 11.</p>

Deflater is a small desktop app that turns off the ads, nags, and junk built
into Windows 11, and keeps them off. It is made to be handed to a friend: one
exe, a clear list of fixes in plain English, and honest labels about what each
one costs you.

## What it does

- **Switches off the noise.** Lock screen ads, File Explorer promos, "finish
  setting up your device" screens, suggestion popups, Bing in Start search,
  search box doodles, Widgets, Edge nags, and the pipeline that installs
  promoted apps on its own.
- **Removes the junk.** Preinstalled apps nobody asked for, each one an
  individual choice, every one reinstallable from the Microsoft Store. Apps
  that can no longer be reinstalled are deliberately not offered.
- **Guards the door.** Blocks manufacturers from auto-installing companion
  apps when you plug in hardware (the switch behind the recent LG incident),
  and can watch for apps that appear without you asking and notify you with a
  one-click remove.
- **Turns AI features off until you want them.** Copilot, Recall snapshots,
  and Click to Do.
- **Trims what gets sent home.** Advertising ID, tailored experiences,
  activity history, typing personalization, and diagnostic data down to the
  required minimum. Honest note: on Home and Pro that is a minimum, not off.
- **Keeps it fixed.** An optional scheduled task re-checks after sign-in and
  weekly, because Windows updates love to bring things back.

## What it never touches

Defender, Secure Boot, TPM, virtualization security, driver delivery, Windows
Update itself, and anything Xbox or Game Pass. Games with kernel anti-cheat
see a completely stock security posture. This is enforced by a unit test over
the fix catalog, not just by good intentions.

## Profiles

Three starting points on one dial, reviewed and editable before anything runs:

| Profile | What it means |
| --- | --- |
| Light Touch | Switches only. Removes nothing, nothing visibly missing. |
| Clean Sweep | Light Touch plus junk apps gone and Bing out of Start. The default. |
| Full Deflate | Everything removable goes. Reinstall what you miss, free. |

Every fix shows its live status, a one-line summary, and an expandable
explanation: what it changes, what you give up, and how to undo it.

## How it works

The app opens without administrator rights and changes nothing until you hit
Apply, at which point Windows shows its standard permission prompt. Switches
are registry policies and Settings-backed toggles, reverted by restoring the
Windows default. App removals go through PowerShell's supported Store app
commands, with deprovisioning so feature updates do not re-seed them. The
maintenance task runs `Deflater.exe --maintenance` headless as the signed-in
user, re-applies anything that drifted, and diffs the installed app list to
catch silent arrivals.

Config and logs live in `%LOCALAPPDATA%\Deflater`. Everything the app does is
written to the log, and the log folder is one click away in the footer.

## Building

Prerequisites: [Go](https://go.dev) 1.26+, [Node](https://nodejs.org) 20+, and
the [Wails v2 CLI](https://wails.io) (`go install
github.com/wailsapp/wails/v2/cmd/wails@latest`).

```powershell
wails build          # production exe in build/bin/Deflater.exe
wails dev            # live-reload development (also serves http://localhost:34115)
```

Checks, from `frontend/`:

```powershell
npm run lint         # Biome
npm run check        # svelte-check (TypeScript)
npm test             # vitest, selection and diff logic
npm run icons        # regenerate icons from assets/icon.svg
```

And from the repo root, `go test ./...` covers the fix catalog: id integrity,
profile nesting, string coverage, and the forbidden-areas guarantee.

## Project layout

```
app.go, main.go        Wails app bridge and entry point (headless mode lives here)
internal/catalog       The fix catalog: every registry value and package name
internal/engine        Applies, reverts, and reads status
internal/maintain      The --maintenance pass the scheduled task runs
internal/watcher       Silent-install detection
internal/...           Small single-purpose packages: reg, appx, schtask, toast
frontend/src           Svelte 5 UI; all copy in lib/strings/en.ts for future languages
assets/                Logo and icon sources (SVG)
```

## License

MIT. Built for friends; read the list, untick what you use, and share it on.
