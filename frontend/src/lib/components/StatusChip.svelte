<script lang="ts">
  import { S } from "../i18n";
  import type { FixKind, FixStatus } from "../types";

  let {
    status,
    kind,
    selected,
    pending,
  }: { status: FixStatus; kind: FixKind; selected: boolean; pending: boolean } = $props();

  const isApp = $derived(kind === "app-junk" || kind === "app-might");

  // The chip answers, in the row's own vocabulary, either "what is this
  // right now?" or, when the toggle diverges from reality, "what will
  // happen when you apply?" Future tense implies the present, so the two
  // are never shown together.
  const chip = $derived.by((): { label: string; tone: string } => {
    if (pending) {
      if (isApp) return { label: S.status.willRemove, tone: "change" };
      if (kind === "onedrive")
        return { label: selected ? S.status.willBlock : S.status.willUnblock, tone: "change" };
      return { label: selected ? S.status.willTurnOn : S.status.willTurnOff, tone: "change" };
    }
    if (status === "unknown") return { label: S.status.unknown, tone: "neutral" };
    if (isApp) {
      return status === "removed"
        ? { label: S.status.notInstalled, tone: "good" }
        : { label: S.status.installed, tone: "neutral" };
    }
    if (kind === "onedrive") {
      if (status === "on") return { label: S.status.blocked, tone: "good" };
      if (status === "partial") return { label: S.status.partlyBlocked, tone: "mixed" };
      return { label: S.status.notBlocked, tone: "neutral" };
    }
    if (status === "on") return { label: S.status.on, tone: "good" };
    if (status === "partial") return { label: S.status.partlyOn, tone: "mixed" };
    return { label: S.status.off, tone: "neutral" };
  });
</script>

<span class="chip {chip.tone}">{chip.label}</span>

<style>
  .chip {
    flex: none;
    font-size: 11.5px;
    line-height: 1;
    padding: 5px 9px;
    border-radius: var(--r-chip);
    background: var(--bg-raised);
    color: var(--text-dim);
    white-space: nowrap;
  }
  .chip.good {
    background: var(--sage-soft);
    color: var(--sage);
  }
  .chip.mixed {
    background: var(--gold-soft);
    color: var(--gold);
  }
  .chip.change {
    background: var(--coral-soft);
    color: var(--coral-bright);
  }
</style>
