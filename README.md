<p align="center">
  <img src="assets/mascot.png" width="420" alt="The Deflater mascot: four worn vinyl cushions arranged like the Windows logo, one deflated, wrapped in a red DEFLATED banner" />
</p>

<h1 align="center">Deflater</h1>

A small Windows 11 app that turns off ads, nags, and preinstalled junk, and keeps them off.

**Version 0.1. This is an alpha of an alpha.** It is a hobby project built to share with friends and family, it is probably not very good, and it is barely tested. It edits your Windows registry and uninstalls apps. Do not use it if you care about your machine. Use at your own risk, and read what a fix does before applying it.

## What it does

- Turns off ads and suggestions: lock screen tips, File Explorer banners, Start menu recommendations, upsell popups
- Uninstalls preinstalled apps you choose. Everything offered can be reinstalled from the Microsoft Store, and removed apps get a Reinstall button that takes you straight there
- Blocks Windows from silently installing promoted apps, and blocks hardware makers from auto-installing their own software when you plug something in
- Turns off Copilot, Recall, and other AI features
- Reduces what the PC sends back to Microsoft: advertising ID, tailored experiences, diagnostic data down to the minimum Windows allows
- Optional: a background task re-applies your choices after Windows updates, and warns you when an app installs itself without asking

The app opens without administrator rights and changes nothing until you press Apply, at which point Windows shows its normal permission prompt.

## What it never touches

Defender, Secure Boot, TPM, driver updates, Windows Update itself, or anything Xbox / Game Pass. A unit test fails the build if a fix tries.

## Undoing things

Switches turn back off inside the app and restore the Windows defaults. Uninstalled apps come back from the Microsoft Store. Settings and logs live in `%LOCALAPPDATA%\Deflater`, and there is an "Open logs" link in the app footer.

## Building it yourself

You need [Go](https://go.dev) 1.26+, [Node](https://nodejs.org) 20+, and the [Wails v2 CLI](https://wails.io):

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build          # produces build/bin/Deflater.exe
wails dev            # development with live reload
```

Checks:

```powershell
go test ./...        # from the repo root
npm run lint         # from frontend/
npm run check
npm test
```

## Contributing

Issues and pull requests are welcome.

## License

[GPL-3.0](LICENSE). Roughly: use it, change it, share it, sell it if you must, but if you distribute your own version you have to publish its source under this same license and keep the credits. There is no warranty of any kind.
