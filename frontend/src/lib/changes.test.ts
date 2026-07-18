import { describe, expect, it } from "vitest";
import { computeChanges, initialSelection, matchesProfile, profileSelection } from "./changes";
import type { FixState } from "./types";

function fix(partial: Partial<FixState> & { id: string }): FixState {
  return {
    category: "ads-nags",
    kind: "switch",
    caution: false,
    profiles: ["light-touch", "clean-sweep", "full-deflate"],
    status: "off",
    ...partial,
  };
}

describe("computeChanges", () => {
  it("applies selected fixes that are off or partial", () => {
    const fixes = [
      fix({ id: "a", status: "off" }),
      fix({ id: "b", status: "partial" }),
      fix({ id: "c", status: "on" }),
    ];
    const { enable, disable } = computeChanges(fixes, new Set(["a", "b", "c"]));
    expect(enable).toEqual(["a", "b"]);
    expect(disable).toEqual([]);
  });

  it("removes selected apps that are still installed", () => {
    const fixes = [fix({ id: "app", kind: "app-junk", status: "installed" })];
    expect(computeChanges(fixes, new Set(["app"])).enable).toEqual(["app"]);
  });

  it("reverts deselected switches that are on", () => {
    const fixes = [fix({ id: "a", status: "on" }), fix({ id: "b", status: "off" })];
    const { enable, disable } = computeChanges(fixes, new Set());
    expect(enable).toEqual([]);
    expect(disable).toEqual(["a"]);
  });

  it("leaves deselected partial fixes alone instead of half-undoing them", () => {
    const fixes = [fix({ id: "half", status: "partial" })];
    const { enable, disable } = computeChanges(fixes, new Set());
    expect(enable).toEqual([]);
    expect(disable).toEqual([]);
  });

  it("never tries to revert an app removal", () => {
    const fixes = [fix({ id: "app", kind: "app-might", status: "removed" })];
    expect(computeChanges(fixes, new Set()).disable).toEqual([]);
  });

  it("reverts the OneDrive policy when deselected", () => {
    const fixes = [fix({ id: "od", kind: "onedrive", status: "on" })];
    expect(computeChanges(fixes, new Set()).disable).toEqual(["od"]);
  });

  it("is a no-op when selection matches reality", () => {
    const fixes = [
      fix({ id: "a", status: "on" }),
      fix({ id: "app", kind: "app-junk", status: "removed" }),
      fix({ id: "b", status: "off" }),
    ];
    const { enable, disable } = computeChanges(fixes, new Set(["a", "app"]));
    expect(enable).toEqual([]);
    expect(disable).toEqual([]);
  });
});

describe("initialSelection", () => {
  it("selects what is already in effect plus managed drift", () => {
    const fixes = [
      fix({ id: "on-now", status: "on" }),
      fix({ id: "gone-app", kind: "app-junk", status: "removed" }),
      fix({ id: "drifted", status: "off" }),
      fix({ id: "untouched", status: "off" }),
    ];
    const sel = initialSelection(fixes, ["drifted"]);
    expect(sel).toEqual(new Set(["on-now", "gone-app", "drifted"]));
  });

  it("never queues a change at startup, including for partial fixes", () => {
    const fixes = [
      fix({ id: "half", status: "partial" }),
      fix({ id: "done", status: "on" }),
      fix({ id: "untouched", status: "off" }),
    ];
    const sel = initialSelection(fixes, []);
    expect(sel).toEqual(new Set(["done"]));
    expect(computeChanges(fixes, sel)).toEqual({ enable: [], disable: [] });
  });

  it("re-queues a managed fix that drifted to partial", () => {
    const fixes = [fix({ id: "managed-drift", status: "partial" })];
    const sel = initialSelection(fixes, ["managed-drift"]);
    expect(computeChanges(fixes, sel).enable).toEqual(["managed-drift"]);
  });

  it("drops managed ids no longer in the catalog", () => {
    const sel = initialSelection([fix({ id: "a" })], ["retired-fix"]);
    expect(sel).toEqual(new Set());
  });
});

describe("profileSelection", () => {
  const fixes = [
    fix({ id: "mild", profiles: ["light-touch", "clean-sweep", "full-deflate"] }),
    fix({ id: "medium", profiles: ["clean-sweep", "full-deflate"] }),
    fix({ id: "max", kind: "app-might", profiles: ["full-deflate"], status: "installed" }),
    fix({ id: "already-gone", kind: "app-junk", profiles: ["full-deflate"], status: "removed" }),
  ];

  it("selects the profile's fixes", () => {
    expect(profileSelection(fixes, "clean-sweep")).toEqual(
      new Set(["mild", "medium", "already-gone"]),
    );
  });

  it("keeps removed apps selected even outside the profile", () => {
    expect(profileSelection(fixes, "light-touch").has("already-gone")).toBe(true);
  });

  it("matchesProfile identifies the active profile exactly", () => {
    const sel = profileSelection(fixes, "clean-sweep");
    expect(matchesProfile(fixes, sel, "clean-sweep")).toBe(true);
    expect(matchesProfile(fixes, sel, "light-touch")).toBe(false);
    expect(matchesProfile(fixes, sel, "full-deflate")).toBe(false);
  });
});
