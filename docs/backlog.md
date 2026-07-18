# Deflater backlog

Notes captured during testing. Not yet actioned — just tracked.

## Behavior & consistency

- **OneDrive has no Reinstall button** unlike other app removals. Reason:
  OneDrive is the `onedrive` fix kind (a Win32 app removed by its own
  uninstaller + a policy), not a Microsoft Store app, so there's no Store
  product page to link to. Decide how to present this consistently
  (e.g. a "how to get it back" note, or a link to microsoft.com/onedrive).
- **OneDrive fix does two jobs** (sets a policy AND uninstalls). More
  broadly: consider **sub-settings** for fixes that bundle actions, so a
  user can e.g. keep an app installed but turn off just its nags.

  DESIGN PROPOSAL (sub-settings) — needs sign-off before building:
  - Today the only genuinely bundled fix is `app-onedrive`: it writes the
    `DisableFileSyncNGSC` policy AND runs Microsoft's uninstaller. Every
    other fix is already single-purpose, so this is a small surface, not a
    catalog-wide refactor.
  - Recommended approach: NOT a generic nested-toggle engine (over-built
    for one case, and nested toggles complicate the apply/revert/snapshot
    model and the profiles). Instead, express a bundle as two normal
    catalog fixes with a lightweight "parent/child" grouping for display:
    e.g. `onedrive-block` (Switch: the policy, reversible) and
    `onedrive-uninstall` (the uninstaller, reinstall-from-Store). The row
    renders them as one card with a primary toggle and a secondary
    "also uninstall it" sub-toggle, but underneath they're independent
    fixes the engine already knows how to apply and revert.
  - This reuses all existing machinery (status, snapshots, profiles,
    refresh, maintenance) with only a UI grouping field (e.g.
    `Group string` on Fix, plus a `Primary bool`) instead of a new
    apply/revert code path. Profiles pick children independently.
  - Open question for the user: is OneDrive the only case worth splitting,
    or are there others (e.g. an app whose nags could be silenced without
    removing it)? That answer decides whether even this light grouping is
    worth it, or whether a doc note ("uninstall is all-or-nothing") is
    enough for now.

## LG companion app — mechanism & our fix (2026-07-18)

Reproduced live: plugging in an LG UltraGear monitor makes DeviceSetup-
Manager download+install `LGElectronics.LGMonitorApp` (Store product
9PM9N6F47JB8). The request originates from the monitor's **driver
software component** (`oem45.inf`, "LG Monitor Support Application",
`SWD\DRIVERENUM`, class SoftwareComponent), but DSM is what performs the
network download.

Our `device-metadata-off` fix sets exactly the two values the community
fixes recommend, verified against the ADMX on-machine:
- gpedit "Prevent automatic download of applications associated with
  device metadata" = `SOFTWARE\Policies\...\Device Metadata\
  PreventDeviceMetadataFromNetwork` (policy `DeviceMetadata_
  PreventDeviceMetadataFromNetwork` in DeviceSetup.admx).
- Control Panel Hardware tab "No, don't allow installs" = the non-policy
  twin `...\CurrentVersion\Device Metadata\PreventDeviceMetadataFromNetwork`.
So the fix targets the right lever. The policy name literally covers
"download of applications associated with device metadata", which is what
DSM does here.

- **STILL TO CONFIRM empirically:** clean-slate replug with block ON must
  actually stop the download. Positive control (block off, deprovisioned)
  must actually install it. Machine-history caveat: a provisioned copy
  lingers after a user uninstall, so a clean repro needs deprovisioning.
- This is the app's headline use case; keep the empirical proof.

### FINAL TEST CONCLUSION (2026-07-18) — cannot verify any block on this box
We NEVER reproduced a fresh install. The only time the app existed was
the pre-test provisioned copy. In EVERY config since removing it — device-
metadata off/on, Deprovisioned entry present/removed, consumer-features
off, silent-app-installs off, pre/post reboot, old debloat task deleted —
DSM requested (166) + located (167) the app but never deployed it
(elevated checks: not installed/provisioned, no AppX deployment).

Both theories tested and FALSIFIED:
- device-metadata-off: DSM still requests/locates with it ON. Unverified.
- "Deprovisioned entry is the block": DELETED it, app STILL didn't
  install. So that entry was NOT what kept it away. (Earlier note that
  deprovision "provably" blocks it was WRONG — corrected here.)

Real cause unknown but almost certainly a per-user "user removed this,
don't re-push" record (StateRepository / per-account Store state) that we
can't cleanly reset on a live machine. Net: zero positive controls, so no
block is verifiable here.

ONLY definitive test: a machine/VM/user account that has NEVER had this
app, with the candidate block set BEFORE first connect.

**Product decision (user, 2026-07-18): PREVENT, don't uninstall/reinstall.**
Keep `device-metadata-off` as the LG block. Do NOT add an uninstall fix —
the user wants an actual block, not remove-and-let-it-come-back.
- `device-metadata-off` is the documented lever (PreventDeviceMetadata-
  FromNetwork, both policy + CurrentVersion) and is what Home-user reports
  say works on a fresh connect. Ship it as the LG protection.
- Honesty in copy: describe it as "blocks manufacturers auto-installing
  their apps" (which is its documented function). Do NOT over-claim a
  guaranteed LG block until verified on a never-touched machine — one
  negative signal stands (DSM still requested/located with block ON), and
  the download-stage block was unprovable on this history-laden machine.
- Deprovision-as-block was proven to work, but it's the uninstall route
  the user declined; keep it only as a fallback idea, not the primary.
- Future verification: a clean VM / never-connected machine with
  device-metadata-off set BEFORE first connect is the only way to confirm.

### EARLIER TEST NOTES — efficacy of device-metadata-off UNVERIFIED
Controlled test on this machine (LG UltraGear, app fully deprovisioned +
device/SWC nodes removed + metadata cache cleared before each replug):
- Block OFF, clean slate: DSM requested (evt 166) + located (167) the app,
  but it did NOT install (elevated check: not installed/provisioned, no
  AppX deployment).
- Block ON (Policy+CurrentVersion=1), clean slate: IDENTICAL — 166/167
  fired, app did NOT install.
Conclusions:
- `PreventDeviceMetadataFromNetwork` does NOT stop DSM requesting/locating
  the companion app (evt 166/167 fire regardless of block).
- App didn't install in EITHER state, so the block made no observable
  difference. The real blocker in both runs is something else — likely
  (a) Windows suppressing re-install of a just-removed app, and/or (b) an
  already-applied fix (`silent-app-installs`: OemPreInstalledAppsEnabled=0,
  DisableWindowsConsumerFeatures=1) blocking it in both runs.
- We CANNOT currently claim device-metadata-off blocks the LG app.
Next to resolve, in order of rigor:
1. Clean machine / VM / brand-new local user that has never seen the
   monitor, with block set BEFORE first connect. Gold standard.
2. On this machine: revert ALL Deflater fixes + `pnputil /delete-driver
   oem45.inf /uninstall` so the driver+SWC+app come fresh (defeats the
   recent-removal suppression) -> replug = positive control -> then apply
   ONE fix at a time to find which actually blocks it (suspect
   silent-app-installs' OEM/consumer-features values, NOT device-metadata).
3. Note: the app may register the Store app at next LOGON (evt 821 seen in
   the original install), so a sign-out/in may be needed to see the real
   end state — check that too.

### Windows Home matters (majority of users)
- `gpedit.msc` does NOT exist on Windows Home — the group-policy route
  the news articles cite is useless to most users. This is a core reason
  the app exists: it writes the registry directly, which works on Home.
- Per MS docs, the POLICY value
  (`SOFTWARE\Policies\...\Device Metadata\PreventDeviceMetadataFromNetwork`)
  OVERRIDES the Settings-dialog value
  (`...\CurrentVersion\Device Metadata\PreventDeviceMetadataFromNetwork`);
  if the policy is unset, the dialog value governs. Our fix sets BOTH, so
  we're covered on Home and Pro. Multiple Home users on Reddit confirm the
  Settings-dialog value alone works.
- Confirm the policy value is actually honored on Home (a Reddit commenter
  wasn't sure). If Home ignores the Policies key, the CurrentVersion value
  (which we also set) is the one that counts there.

## Testing to do

- Verify **blocking manufacturer auto-installs** actually prevents an
  install in practice (not just that the reg value is set).
- Test the **automated detection (watcher)** end to end.
- Determine **how "block manufacturer auto-installs" and the watcher
  interact** with each other.

## UX / features

- **Search/filter box at the top, sticky** on scroll.
- **Minimal sticky branding** on scroll — try it, see how it looks.
- **In-app toast on apply** to confirm what happened. The single top
  banner is ambiguous (can't tell old vs new).
- **Prefer toasts over the top banner** for events. For results that need
  a restart, use a small modal / message / richer toast that explains the
  outcome.
- **Reword the elevation warning:** "standard Windows prompt" → call it a
  UAC prompt, or something less technical, and **explain why the app has
  to restart** in plain language for non-technical users.

## Windows edition awareness

- DONE (detection + display): `internal/winver` reads the edition; the
  report carries `edition`/`home` and the footer shows e.g. "Windows 11
  Home". STILL TO DO: the per-fix tailoring below (hiding/marking genuine
  Home no-ops) needs a careful per-fix audit before shipping, so it's not
  wired yet — mislabeling a working fix as a no-op would be worse than
  saying nothing. `silent-app-installs` is the tricky case: its
  CloudContent policies are Enterprise/Edu-only, but it still works on
  Home via the HKCU ContentDeliveryManager values, so it is NOT a no-op.
- For the device-metadata / LG fix, no edition detection is needed: set
  both the policy value and the CurrentVersion value always (covers Home
  and Pro; policy overrides the dialog where honored). gpedit == the
  policy registry key, so writing the key IS "setting the group policy".
- Where edition detection WOULD help (separate enhancement): hide or mark
  fixes that are genuine no-ops on Home (e.g. `DisableWindowsConsumer-
  Features` is Enterprise/Education-only), tailor messaging ("works on
  your edition"), and adjust anything with edition-specific behavior
  (telemetry floor, some CloudContent policies). Detect via
  `Get-CimInstance Win32_OperatingSystem` Caption / registry EditionID.

## Visual / layout

- **More space between sections.** They're compact and run together; add
  negative space so sections are easier to tell apart.

## Verified working (testing log)

- All **Reinstall buttons** open the correct Microsoft Store product page
  (tested across every app removal).

## Conflicting automation from other tools

- **Detect and offer to remove the old win11-debloat scheduled task.**
  Deflater is the successor to the `win11-debloat` script. A machine that
  ran the old project has a `Win11-Debloat-Maintenance` task that re-runs
  `debloat-win11.ps1` at every logon, silently re-applying the old
  script's settings — which overrides/conflicts with Deflater's state
  (confirmed live 2026-07-18: it undid a manual change after reboot).
  Deflater should detect this task (and ideally other known debloat
  tools' tasks) and offer to remove it, since it now owns this job.
- Testing implication: this old task was quietly re-applying
  `silent-app-installs`-style values during the LG tests, contaminating
  them. Disable it before any clean LG repro.

## Update awareness (no update service)

- **Check for a newer Deflater release and link to it.** On launch (or
  once a day), compare `appVersion` against the latest GitHub release; if
  newer, show a small unobtrusive "new version available" note linking to
  the releases page (github.com/Sigler/deflater/releases). NOT an
  auto-updater — just awareness + a link. Cheap: hit the GitHub releases
  API or a static version file; fail silent if offline.

## Live state / reactivity

- **Detect reinstalls while the app is open.** After reinstalling an app
  from the Store, its Reinstall button stays stale until a manual reload.
  Add a lightweight, non-obnoxious way to notice when a removed app comes
  back while open and update the row (e.g. a gentle heartbeat re-check,
  and/or a rescan on window focus). Find an efficient, elegant approach.
- **Detect external setting changes while open.** If another process or
  app changes one of these settings while Deflater is open, notice it,
  tell the user, and update the toggle/button state to match reality
  (registry change notifications, or the same heartbeat).

## Build provenance & token spend

- **Model:** Claude Fable 5 (`claude-fable-5`) for the entire build thread.
- **Effort/reasoning level:** not reliably self-reportable from inside the
  session; check the Claude Code config for the exact setting.
- **Sub-agent spend (measured, from agent completion reports):** ~518,518
  output tokens across 10 background agents — catalog verification (63,537),
  IA/profiles (20,497), UX research (34,658), Store IDs (30,052), and six
  review agents (security 62,163, build/release 42,181, concurrency 59,111,
  catalog 34,540, error-handling 81,759, frontend 90,020).
- **Whole-thread total:** NOT available from inside the session — no tool
  exposes a running token count. Get the true total from the Claude Code /
  claude.ai usage view for this session. Going forward, per-action token
  spend also can't be self-measured; track it from that usage view.

## Apply & refresh (avoiding restarts)

Some fixes only take effect after the process that reads them restarts.
Confirmed live: **websearch-off** (`DisableSearchBoxSuggestions`) needs
only a **SearchHost/Explorer restart** — NOT a full sign-out or reboot —
before Start-menu web results disappear. This is the model to chase.

- **Auto-refresh after apply.** When a fix that needs it is applied,
  Deflater could restart just the relevant shell process (SearchHost,
  Explorer, etc.) so the change takes effect immediately, turning "sign
  out to see this" into instant. Do it gently and only for fixes that
  need it; warn the user first if it's disruptive (Explorer flicker).
- **Tell the user what a fix needs.** After applying, show (toast/modal)
  whether a fix is live now, needs a process refresh (offer to do it),
  needs sign-out, or needs a full restart.
- **Research + test the whole catalog for refresh requirements.** For
  each fix determine the *minimum* refresh needed: nothing / Explorer /
  SearchHost / sign-out / reboot, and find the lightest option that
  works. Goal: avoid a full system restart wherever possible.
- Current copy already says some fixes "take effect after your next
  sign-in" — replace that with the specific, lightest action per fix.

## Recall

- **Surface existing Recall snapshots without toggling anything.** Detect
  whether Recall snapshots exist on the machine and show it up front: that
  they exist, how much disk space they use, where they live. If feasible,
  point the user to how to back them up or remove them. It's a privacy
  signal worth making visible before the user touches the Recall fixes.
