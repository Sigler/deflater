<script lang="ts">
  import { S } from "../i18n";
  import type { FixKind, FixStatus } from "../types";

  let { status, kind }: { status: FixStatus; kind: FixKind } = $props();

  const isApp = $derived(kind === "app-junk" || kind === "app-might");

  const label = $derived.by(() => {
    if (status === "removed") return S.status.notInstalled;
    if (status === "installed") return S.status.installed;
    return S.status[status] ?? S.status.unknown;
  });

  // For switches, "on" is the good state. For apps, gone is the goal but
  // "installed" is simply neutral fact, not an error.
  const tone = $derived.by(() => {
    if (status === "on" || status === "removed") return "good";
    if (status === "partial") return "mixed";
    return "neutral";
  });
</script>

<span class="chip {tone}" class:app={isApp}>{label}</span>

<style>
  .chip {
    flex: none;
    font-size: 11.5px;
    line-height: 1;
    padding: 5px 9px;
    border-radius: var(--r-chip);
    background: var(--bg-raised);
    color: var(--text-dim);
    border: 1px solid transparent;
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
