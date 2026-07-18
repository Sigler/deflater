// Every user-facing word in Deflater lives here, keyed by the ids the
// Go catalog exposes. Adding a language means translating this file and
// switching the export in ../i18n.ts. Style notes: plain sentences,
// honest about tradeoffs, no jargon, no drama.

export interface FixStrings {
  title: string;
  summary: string;
  what: string;
  tradeoff?: string;
  undo: string;
  // Microsoft Store search term for uninstalled apps' Reinstall link.
  store?: string;
}

export const en = {
  app: {
    name: "Deflater",
    tagline: "Make Windows slightly less obnoxious.",
    scanChecked: "Checked your current Windows settings",
    scanApplied: (n: number, total: number) => `${n} of ${total} applied`,
    loading: "Reading your system…",
    loadingHint: "Checking every switch and installed app. This takes a few seconds.",
    loadFailed: "Couldn't read your system.",
    retry: "Try again",
  },

  search: {
    placeholder: "Search settings",
    clear: "Clear search",
    noResults: (q: string) => `Nothing matches "${q}".`,
  },

  profiles: {
    label: "Preset",
    custom: "Custom selection",
    selected: (n: number, total: number) => `${n} of ${total} selected`,
    "light-touch": {
      title: "Light Touch",
      tagline: "Ads, nags, and tracking off. Nothing removed, nothing you would notice missing.",
    },
    "clean-sweep": {
      title: "Clean Sweep",
      tagline: "Light Touch, plus Bing out of your Start menu and the junk apps gone.",
      badge: "Recommended",
    },
    "full-deflate": {
      title: "Full Deflate",
      tagline: "Everything removable goes. Reinstall anything you miss from the Store, free.",
    },
  },

  categories: {
    "ads-nags": {
      title: "Ads and nags",
      nav: "Ads and nags",
      blurb: "Promotions, upsells, and pestering built into Windows itself.",
    },
    "junk-apps": {
      title: "Junk apps",
      nav: "Junk apps",
      blurb: "Preinstalled apps nobody asked for, and the pipes that sneak new ones in.",
    },
    "start-search": {
      title: "Start menu, search, and taskbar",
      nav: "Start and search",
      blurb: "Clutter in the places you use most.",
    },
    "copilot-ai": {
      title: "Copilot and AI",
      nav: "Copilot and AI",
      blurb: "Windows AI features, off until you actually want them.",
    },
    privacy: {
      title: "Privacy",
      nav: "Privacy",
      blurb: "What this PC quietly sends back to Microsoft.",
    },
    "might-use": {
      title: "Apps you might use",
      nav: "Apps you might use",
      blurb:
        "Preinstalled, but some people use these on purpose. Untick anything you actually use.",
    },
  },

  nav: {
    maintenance: "Maintenance",
  },

  status: {
    applied: "Applied",
    notApplied: "Not applied",
    willApply: "Will apply",
    willUninstall: "Will uninstall",
    willUndo: "Will undo",
    unknown: "Unknown",
  },

  badges: {
    willChange: "Changes when you apply",
  },

  details: {
    what: "What it does",
    tradeoff: "What you give up",
    undo: "Undoing it",
    mechanismReg: (n: number) =>
      n === 1 ? "Sets 1 Windows setting." : `Sets ${n} Windows settings.`,
    mechanismApp: (names: string) => `Uninstalls: ${names}.`,
    undoSwitch: "Flip the toggle off and apply. Deflater restores the Windows default.",
    undoApp: "Reinstall it from the Microsoft Store any time, free.",
    reinstall: "Reinstall",
    reinstallHint: "Opens the Microsoft Store",
    partialNote:
      "Part of this fix is already set on this PC, most likely by another tool or an older script. Applying it completes the rest.",
  },

  apply: {
    changesReady: (n: number) =>
      n === 1
        ? "1 change pending. Nothing is applied yet."
        : `${n} changes pending. Nothing is applied yet.`,
    applyCount: (n: number) => (n === 1 ? "Apply 1 change" : `Apply ${n} changes`),
    apply: "Apply changes",
    applying: "Applying…",
    reset: "Reset",
    doneTitle: "Done",
    doneBody:
      "All changes applied. A few finish after you sign out and back in, so do that when convenient.",
    doneMaintenanceTip:
      "Tip: turn on 'Keep it fixed automatically' near the bottom, so a Windows update cannot quietly undo this.",
    doneSomeFailed: (n: number) =>
      n === 1 ? "1 change didn't apply." : `${n} changes didn't apply.`,
    saveWarning:
      "Your changes were applied, but Deflater couldn't save a record of them. Maintenance may not track them. Details are in the logs.",
    applyError: "Something went wrong applying your changes:",
    elevateTitle: "Windows will ask for permission",
    elevateBody:
      "Changing these settings needs administrator rights, and Windows only grants those to a program as it starts. So Deflater has to close and reopen once, which is why you'll see Windows' blue User Account Control (UAC) box asking you to allow it. Say yes, and Deflater picks up exactly where you left off.",
    elevateConfirm: "Continue",
    elevateCancel: "Not now",
    resuming: "Continuing your changes…",
    doneClean: "All changes applied.",
    refreshNone: "Everything is live now.",
    refreshExplorer: "One change needs an Explorer restart to show, which flickers the taskbar for a second.",
    refreshSignout: "A few changes finish after you sign out and back in.",
    refreshReboot: "A few changes finish after you restart your PC.",
    restartExplorer: "Restart Explorer now",
    restartingExplorer: "Restarting…",
  },

  toast: {
    dismiss: "Dismiss",
  },

  maintenance: {
    sectionBlurb: "Keep your choices applied, and get warned when apps sneak in.",
    mismatch:
      "Automatic maintenance is set to on, but its scheduled task isn't registered. Apply any change and allow the User Account Control (UAC) prompt to set it up.",
    title: "Keep it fixed automatically",
    body: "Windows updates love to bring junk back. Deflater can quietly re-check after every sign-in and once a week, and re-apply your choices when something drifts.",
    on: "On",
    off: "Off",
    pendingElevation: "Turns on when you next apply changes.",
    watcherTitle: "Warn me when apps install themselves",
    watcherBody:
      "Checks after sign-in and weekly. If an app appears that you did not install, like a manufacturer app arriving with a new device, you get a notification and can remove it here. Works on its own, with or without the switch above.",
  },

  alerts: {
    title: "Installed without asking",
    body: "These appeared since Deflater last looked, and you did not install them through the usual ways.",
    remove: "Remove",
    removing: "Removing…",
    dismiss: "Dismiss all",
  },

  conflicts: {
    title: "Another cleanup tool is running here",
    body: "Deflater found a leftover scheduled task from an earlier debloat tool. It re-applies its own settings on a schedule, which can quietly fight your choices here. Deflater does this job now, so the old task is safe to remove.",
    remove: "Remove it",
    removing: "Removing…",
    // {tool} names the tool the task belongs to.
    fromTool: (tool: string) => `from ${tool}`,
    removed: "Removed the conflicting scheduled task.",
  },

  footer: {
    logs: "Open logs",
    assurance: "Never touches Defender, Secure Boot, TPM, or anything Xbox or Game Pass.",
    version: (v: string) => `Deflater ${v} · anti-bloat enforcement, est. 2026`,
    // {v} is the newer version available on GitHub.
    updateAvailable: (v: string) => `Version ${v} is available`,
  },

  fixes: {
    // ---- Ads and nags ---------------------------------------------------
    "lockscreen-ads": {
      title: "Block lock screen ads",
      summary: "Stops the tips, quizzes, and app offers on your lock screen.",
      what: "Turns off the two switches behind the lock screen's rotating 'fun facts, tips, and more', which is where Microsoft slots promotions. Your lock screen wallpaper, including Spotlight photos, stays exactly as you set it.",
      undo: "Flip the toggle off and apply. The Settings switches go back to their defaults.",
    },
    "explorer-ads": {
      title: "Hide File Explorer ads",
      summary: "Removes OneDrive and Office promo banners inside File Explorer.",
      what: "Turns off 'sync provider notifications', the mechanism that puts subscription offers above your own files. Explorer applies it after your next sign-in.",
      undo: "Flip the toggle off and apply.",
    },
    "scoobe-off": {
      title: "Turn off 'finish setting up' screens",
      summary: "Stops the full-screen 'Let's finish setting up your device' interruptions.",
      what: "Turns off the suggestion screen Windows shows after updates that pushes Microsoft 365, OneDrive backup, and a Microsoft account. One per-user setting, the same one buried in Settings under Notifications.",
      undo: "Flip the toggle off and apply.",
    },
    "suggested-toasts-off": {
      title: "Silence 'Get even more' popups",
      summary: "Silences the suggestion notifications in the corner of your screen.",
      what: "Turns off the notification sender Windows uses for upsell toasts. Real notifications from your apps are untouched.",
      undo: "Flip the toggle off and apply.",
    },
    "settings-suggestions": {
      title: "Turn off Start and Settings suggestions",
      summary: "Clears promoted apps and 'tips' from Start and the Settings app.",
      what: "Turns off nine suggestion switches: promoted content in Settings, app suggestions in Start, 'recommendations for tips and new apps', and the welcome experience after updates.",
      tradeoff:
        "Start and Settings stop recommending new apps and features. Most people call that the point.",
      undo: "Flip the toggle off and apply.",
    },
    "edge-nags": {
      title: "Stop Edge nags and background running",
      summary: "Stops Edge running in the background and pushing its sidebar and offers.",
      what: "Sets five Edge policies: no startup boost, no background mode, no Copilot sidebar, no recommendation popups, and no sending your browsing personalization data.",
      tradeoff:
        "Edge opens a moment slower and will show 'Managed by your organization' in its menu. That label is how Edge reports any policy, including these.",
      undo: "Flip the toggle off and apply. The policies are deleted and the label disappears.",
    },
    "edge-shortcut": {
      title: "Stop Edge shortcut resurrection",
      summary: "Stops Edge re-creating its desktop shortcut after every update.",
      what: "Sets the two Edge updater policies that control shortcut creation. Note: applying this also removes an Edge desktop shortcut if one exists now.",
      tradeoff: "If you like having the Edge icon on your desktop, skip this one.",
      undo: "Flip the toggle off and apply, then re-create the shortcut from the Start menu if you want it back.",
    },

    // ---- Junk apps ------------------------------------------------------
    "silent-app-installs": {
      title: "Block self-installing apps",
      summary: "Blocks the pipeline that installs promoted apps without asking.",
      what: "Turns off the 'content delivery' switches Windows uses to silently install suggested Store apps and pin promo tiles, the same pipeline that famously delivered Candy Crush. Also sets the matching machine-wide policies for good measure.",
      undo: "Flip the toggle off and apply.",
    },
    "device-metadata-off": {
      title: "Block manufacturer auto-installs",
      summary: "Blocks Windows from auto-downloading hardware makers' companion apps.",
      what: "Turns off Windows' automatic download of device metadata and the manufacturer apps that ride along with it. This is the documented switch for that behavior, the kind reported in recent cases like LG's monitor app installing itself. Deflater sets it both ways, as the official policy and the Settings-screen value, so it applies on Home and Pro alike. Driver installation through Windows Update is not affected.",
      tradeoff:
        "If a gadget needs a companion app, you install it yourself from the maker's site or the Store.",
      undo: "Flip the toggle off and apply.",
    },
    "app-officehub": {
      title: "Uninstall Microsoft 365 promo hub",
      summary: "Removes the Microsoft 365 app that mostly sells subscriptions.",
      what: "Uninstalls the Microsoft 365 hub app. This is the promo shell, not Office itself: Word, Excel, and your documents are untouched whether or not you have Office installed.",
      undo: "Reinstall 'Microsoft 365' from the Microsoft Store any time.",
      store: "Microsoft 365",
    },
    "app-news": {
      title: "Uninstall News app",
      summary: "Removes Microsoft's Bing-powered news app.",
      what: "Uninstalls Microsoft News for every account on this PC and stops Windows re-adding it for new accounts.",
      undo: "Reinstall 'Microsoft News' from the Microsoft Store any time.",
      store: "Microsoft News",
    },
    "app-weather": {
      title: "Uninstall Weather app",
      summary: "Removes the MSN Weather app.",
      what: "Uninstalls MSN Weather. The taskbar weather button is separate; that belongs to Widgets, covered below.",
      tradeoff: "If you actually open the Weather app, keep it or reinstall it later.",
      undo: "Reinstall 'MSN Weather' from the Microsoft Store any time.",
      store: "MSN Weather",
    },
    "app-solitaire": {
      title: "Uninstall Solitaire",
      summary: "Removes the ad-stuffed Solitaire collection.",
      what: "Uninstalls Microsoft Solitaire Collection, which interrupts card games with video ads unless you pay monthly.",
      tradeoff:
        "If someone in the house plays it daily, keep it. It reinstalls free from the Store.",
      undo: "Reinstall 'Microsoft Solitaire Collection' from the Microsoft Store any time.",
      store: "Microsoft Solitaire Collection",
    },
    "app-gethelp": {
      title: "Uninstall Get Help app",
      summary: "Removes Microsoft's support-chat app.",
      what: "Uninstalls Get Help, the app that opens when Windows wants to route you to Microsoft support articles and chat.",
      tradeoff: "If you ever contact Microsoft support, you would reinstall it first.",
      undo: "Reinstall 'Get Help' from the Microsoft Store any time.",
      store: "Get Help",
    },
    "app-feedback": {
      title: "Uninstall Feedback Hub",
      summary: "Removes the app for sending feedback to Microsoft.",
      what: "Uninstalls Feedback Hub, which most people open exactly once, by accident, from a keyboard shortcut.",
      undo: "Reinstall 'Feedback Hub' from the Microsoft Store any time.",
      store: "Feedback Hub",
    },
    "app-bingsearch": {
      title: "Uninstall Bing Search app",
      summary: "Removes the 'Web Search from Microsoft Bing' component.",
      what: "Uninstalls the Bing web search app added in recent Windows versions. Pairs well with the Start menu web results switch below.",
      undo: "Reinstall it from the Microsoft Store any time.",
      store: "Bing Search",
    },
    "app-powerautomate": {
      title: "Uninstall Power Automate",
      summary: "Removes Microsoft's workflow automation stub.",
      what: "Uninstalls Power Automate Desktop, preinstalled for enterprise automation almost no home user touches.",
      undo: "Reinstall 'Power Automate' from the Microsoft Store any time.",
      store: "Power Automate",
    },

    // ---- Start, search, taskbar ----------------------------------------
    "websearch-off": {
      title: "Turn off Bing in Start search",
      summary: "Start search shows your apps and files, not web results.",
      what: "Sets the policy that removes web suggestions from Start menu search. Takes effect after you sign out and back in.",
      tradeoff:
        "Typing a web question into Start stops showing web answers; use your browser for those. If search ever misbehaves after a Windows update, undoing this is one click.",
      undo: "Flip the toggle off and apply, then sign out and in.",
    },
    widgets: {
      title: "Turn off Widgets and taskbar news",
      summary: "Removes the Widgets board and its taskbar weather-and-news button.",
      what: "Turns off the Widgets feature by policy and hides its taskbar button. The feed, the weather button, and the hover-over news panel all go.",
      tradeoff: "If you use the weather glance or the news feed, skip this one.",
      undo: "Flip the toggle off and apply. The button returns after a sign-out.",
    },
    "search-highlights": {
      title: "Turn off search box doodles",
      summary: "Stops the daily Bing artwork and trending content in the search box.",
      what: "Turns off 'search highlights', the rotating illustrations and trending searches Microsoft pipes into the taskbar search box.",
      undo: "Flip the toggle off and apply.",
    },

    // ---- Copilot and AI -------------------------------------------------
    "app-copilot": {
      title: "Uninstall Copilot",
      summary: "Uninstalls the Copilot app and its taskbar presence.",
      what: "Uninstalls the Copilot app. On current Windows versions this is the mechanism that actually works; the old 'turn off Copilot' policy is deprecated and ignored. The silent-install block above keeps promotions from re-adding it.",
      tradeoff: "If you use Copilot, keep it. It reinstalls free from the Store.",
      undo: "Reinstall 'Microsoft Copilot' from the Microsoft Store any time.",
      store: "Microsoft Copilot",
    },
    "recall-off": {
      title: "Block Recall screen snapshots",
      summary: "Stops Windows taking searchable screenshots of what you do.",
      what: "Turns off Recall's snapshot saving. Fully reversible, and it keeps any snapshots already on the PC. On PCs without Recall it simply locks the door in advance.",
      undo: "Flip the toggle off and apply. Recall can save snapshots again.",
    },
    "recall-purge": {
      title: "Remove Recall and wipe its snapshots",
      summary: "Removes the Recall component and deletes any snapshots already saved.",
      what: "Goes further than the switch above: it removes the Recall feature entirely and permanently deletes every snapshot already stored on this PC. A restart completes it.",
      tradeoff:
        "This deletes data. If you deliberately use Recall to find things you have seen, its entire memory is erased and cannot be recovered.",
      undo: "Flip the toggle off and apply to allow Recall again, but snapshots already deleted are gone for good.",
    },
    "click-to-do-off": {
      title: "Turn off Click to Do",
      summary: "Turns off the AI actions layer over your screen.",
      what: "Sets the policy that disables Click to Do, the feature that analyzes what is on screen to offer AI actions. Harmless no-op on PCs that do not have it.",
      undo: "Flip the toggle off and apply.",
    },

    // ---- Privacy --------------------------------------------------------
    "advertising-id": {
      title: "Stop personalized ad tracking",
      summary: "Stops apps using your advertising ID to profile you.",
      what: "Disables the advertising ID by policy and flips the matching Settings switch. Apps can no longer tie ads to your identity across apps.",
      tradeoff: "You see the same number of ads, just not tailored to you.",
      undo: "Flip the toggle off and apply.",
    },
    "telemetry-minimum": {
      title: "Minimize diagnostic data",
      summary: "Dials Windows diagnostic reporting down to the required minimum.",
      what: "Sets diagnostic data to 'Required', the lowest level Windows Home and Pro honor, and stops Windows asking you for feedback. Honest note: this is a minimum, not off; only Enterprise editions can go lower.",
      undo: "Flip the toggle off and apply.",
    },
    "tailored-experiences": {
      title: "Turn off personalized tips and offers",
      summary: "Stops Microsoft using your diagnostic data to target tips and ads at you.",
      what: "Turns off 'tailored experiences' both as the user policy and the Settings switch, so diagnostic data cannot feed personalized promotions.",
      undo: "Flip the toggle off and apply.",
    },
    "activity-history": {
      title: "Stop recording activity history",
      summary: "Stops Windows recording a timeline of what you open and do.",
      what: "Turns off activity publishing by policy. Windows stops accumulating the activity feed some features use for 'pick up where you left off' suggestions.",
      undo: "Flip the toggle off and apply.",
    },
    "inking-personalization": {
      title: "Turn off typing personalization",
      summary: "Stops Windows building a profile from what you type and write.",
      what: "Turns off inking and typing personalization, the custom dictionary Windows builds from your writing, including the contacts harvesting switch. Autocorrect and normal typing suggestions still work.",
      undo: "Flip the toggle off and apply.",
    },
    "delivery-optimization": {
      title: "Stop uploading updates to strangers",
      summary: "Stops your bandwidth being used to send Windows updates to other people.",
      what: "Limits update sharing to your own local network. Updates download and install exactly as before; your PC just stops seeding them to the internet.",
      undo: "Flip the toggle off and apply.",
    },

    // ---- Apps you might use --------------------------------------------
    "app-onedrive": {
      title: "Uninstall OneDrive",
      summary: "Uninstalls the OneDrive sync app and stops its sign-in nags.",
      what: "Runs Microsoft's own OneDrive uninstaller and sets the policy that keeps it from running. Files already on this PC stay where they are, and everything in the cloud stays at onedrive.com. Nothing is deleted.",
      tradeoff:
        "Syncing and cloud backup stop. If you rely on OneDrive for backup or shared folders, keep it. Files stored online-only need downloading from onedrive.com first.",
      undo: "Flip the toggle off and apply to lift the block, then reinstall OneDrive from microsoft.com/onedrive.",
    },
    "app-phonelink": {
      title: "Uninstall Phone Link",
      summary: "Removes phone integration: texts, calls, and photos on this PC.",
      what: "Uninstalls Phone Link and its cross-device helper together; removing only one leaves phone integration half-broken, so Deflater treats them as a set.",
      tradeoff: "No more phone notifications, texts, or photos on this PC until reinstalled.",
      undo: "Reinstall 'Phone Link' from the Microsoft Store; the helper comes back with it.",
      store: "Phone Link",
    },
    "app-teams": {
      title: "Uninstall Microsoft Teams",
      summary: "Removes the preinstalled Teams app.",
      what: "Uninstalls the Teams app Windows preinstalls and pins. If work or school uses Teams in the browser, that still works.",
      tradeoff: "Keep it if family, work, or school call you on Teams.",
      undo: "Reinstall 'Microsoft Teams' from the Microsoft Store any time.",
      store: "Microsoft Teams",
    },
    "app-outlook": {
      title: "Uninstall new Outlook",
      summary: "Removes the new Outlook mail app.",
      what: "Uninstalls the 'new Outlook for Windows' that replaced Mail and Calendar. Your mail lives on the server and is untouched; other mail apps keep working.",
      tradeoff: "Keep it if it is your mail app.",
      undo: "Reinstall 'Outlook for Windows' from the Microsoft Store any time.",
      store: "Outlook for Windows",
    },
    "app-clipchamp": {
      title: "Uninstall Clipchamp",
      summary: "Removes Microsoft's video editor.",
      what: "Uninstalls Clipchamp.",
      tradeoff: "Keep it if you edit videos with it.",
      undo: "Reinstall 'Clipchamp' from the Microsoft Store any time.",
      store: "Clipchamp",
    },
    "app-todo": {
      title: "Uninstall Microsoft To Do",
      summary: "Removes the To Do list app.",
      what: "Uninstalls Microsoft To Do. Your lists live in your Microsoft account and reappear if you ever reinstall.",
      tradeoff: "Keep it if you use it for lists and reminders.",
      undo: "Reinstall 'Microsoft To Do' from the Microsoft Store any time.",
      store: "Microsoft To Do",
    },
    "app-family": {
      title: "Uninstall Family Safety",
      summary: "Removes the parental controls app.",
      what: "Uninstalls Microsoft Family Safety.",
      tradeoff: "Do not remove this on a child's PC or a PC managed with family screen-time rules.",
      undo: "Reinstall 'Microsoft Family Safety' from the Microsoft Store any time.",
      store: "Microsoft Family Safety",
    },
    "app-quickassist": {
      title: "Uninstall Quick Assist",
      summary: "Removes the remote assistance app scammers love.",
      what: "Uninstalls Quick Assist, the built-in remote control tool. Phone scammers talk victims into opening it far more often than family IT does. Removing it is a small safety upgrade for most homes.",
      tradeoff:
        "If someone you trust helps you remotely with Quick Assist, keep it, or reinstall it together when needed.",
      undo: "Reinstall 'Quick Assist' from the Microsoft Store any time.",
      store: "Quick Assist",
    },
  } as Record<string, FixStrings>,
};
