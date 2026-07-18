// Pure selection-to-changes logic, kept free of UI and bridge concerns
// so it can be unit tested.

import type { FixState } from "./types";

export interface Changes {
  enable: string[];
  disable: string[];
}

// A fix needs applying when the user wants it but the machine disagrees.
function needsApply(fix: FixState, selected: boolean): boolean {
  if (!selected) return false;
  return fix.status === "off" || fix.status === "partial" || fix.status === "installed";
}

// A fix needs reverting when deselected but still in effect. App
// removals cannot be reverted here (the Store reinstalls them), so only
// switches and OneDrive's policy half qualify.
function needsRevert(fix: FixState, selected: boolean): boolean {
  if (selected) return false;
  if (fix.kind === "app-junk" || fix.kind === "app-might") return false;
  return fix.status === "on" || fix.status === "partial";
}

// computeChanges turns the current selection into the enable/disable
// lists the backend applies.
export function computeChanges(fixes: FixState[], selection: Set<string>): Changes {
  const enable: string[] = [];
  const disable: string[] = [];
  for (const fix of fixes) {
    const selected = selection.has(fix.id);
    if (needsApply(fix, selected)) enable.push(fix.id);
    else if (needsRevert(fix, selected)) disable.push(fix.id);
  }
  return { enable, disable };
}

// initialSelection reflects reality at startup: everything already in
// effect (fully or partly, so a partial fix completes rather than
// reverts), plus anything the user manages that has drifted off.
export function initialSelection(fixes: FixState[], managed: string[]): Set<string> {
  const sel = new Set<string>();
  for (const fix of fixes) {
    if (fix.status === "on" || fix.status === "partial" || fix.status === "removed")
      sel.add(fix.id);
  }
  for (const id of managed) sel.add(id);
  // Managed ids may reference fixes retired from the catalog; drop them.
  const known = new Set(fixes.map((f) => f.id));
  for (const id of sel) if (!known.has(id)) sel.delete(id);
  return sel;
}

// profileSelection returns the selection a profile represents, keeping
// already-removed apps selected (they cannot be un-removed here, and
// deselecting them would only confuse the diff).
export function profileSelection(fixes: FixState[], profile: string): Set<string> {
  const sel = new Set<string>();
  for (const fix of fixes) {
    if (fix.profiles.includes(profile)) sel.add(fix.id);
    else if (fix.status === "removed") sel.add(fix.id);
  }
  return sel;
}

// matchesProfile reports whether the selection is exactly what the
// profile would choose, used to highlight the active profile card.
export function matchesProfile(
  fixes: FixState[],
  selection: Set<string>,
  profile: string,
): boolean {
  const expected = profileSelection(fixes, profile);
  if (expected.size !== selection.size) return false;
  for (const id of expected) if (!selection.has(id)) return false;
  return true;
}
