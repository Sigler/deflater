<script lang="ts">
  import { S } from "../i18n";
  import type { FixStatus } from "../types";

  let {
    status,
    selected,
    pending,
  }: { status: FixStatus; selected: boolean; pending: boolean } = $props();

  // One vocabulary for every row: the chip describes the fix the title
  // names, so it can never be misread as the state of the annoyance
  // ("On" next to "Turn off widgets" could; "Applied" cannot). Present
  // tense states reality; coral future tense appears when the toggle
  // diverges and previews what Apply will do.
  const chip = $derived.by((): { label: string; tone: string } => {
    if (pending) {
      return selected
        ? { label: S.status.willApply, tone: "change" }
        : { label: S.status.willUndo, tone: "change" };
    }
    switch (status) {
      case "on":
      case "removed":
        return { label: S.status.applied, tone: "good" };
      case "partial":
        return { label: S.status.partlyApplied, tone: "mixed" };
      case "unknown":
        return { label: S.status.unknown, tone: "neutral" };
      default: // "off" | "installed"
        return { label: S.status.notApplied, tone: "neutral" };
    }
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
