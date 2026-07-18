<script lang="ts">
  import { S } from "../i18n";
  import type { FixKind, FixStatus } from "../types";

  let { status, kind }: { status: FixStatus; kind: FixKind } = $props();

  const isApp = $derived(kind === "app-junk" || kind === "app-might");

  // Chips speak only when there is something worth saying: the fix is
  // already in place, partly in place, or the app is already gone.
  // Untouched rows stay quiet.
  const chip = $derived.by((): { label: string; tone: string } | null => {
    if (status === "on") return { label: S.status.fixed, tone: "good" };
    if (status === "partial") return { label: S.status.partly, tone: "mixed" };
    if (status === "removed" && isApp) return { label: S.status.notInstalled, tone: "good" };
    if (status === "unknown") return { label: S.status.unknown, tone: "neutral" };
    return null;
  });
</script>

{#if chip}
  <span class="chip {chip.tone}">{chip.label}</span>
{/if}

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
</style>
