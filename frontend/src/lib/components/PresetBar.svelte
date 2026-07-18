<script lang="ts">
  import { matchesProfile } from "../changes";
  import { S } from "../i18n";
  import type { FixState } from "../types";

  let {
    fixes,
    selection,
    onpick,
  }: {
    fixes: FixState[];
    selection: Set<string>;
    onpick: (profile: string) => void;
  } = $props();

  const presets = ["light-touch", "clean-sweep", "full-deflate"] as const;

  const active = $derived(presets.find((p) => matchesProfile(fixes, selection, p)) ?? null);
  const hint = $derived(active === null ? S.profiles.custom : S.profiles[active].tagline);
  const selectedCount = $derived(fixes.filter((f) => selection.has(f.id)).length);
</script>

<div class="bar">
  <span class="label" id="preset-label">{S.profiles.label}</span>
  <div class="segments" role="radiogroup" aria-labelledby="preset-label">
    {#each presets as p (p)}
      <button
        type="button"
        role="radio"
        aria-checked={active === p}
        class="segment"
        class:active={active === p}
        title={S.profiles[p].tagline}
        onclick={() => onpick(p)}
      >
        {S.profiles[p].title}
        {#if p === "clean-sweep"}
          <span class="dot" title={S.profiles["clean-sweep"].badge} aria-hidden="true"></span>
        {/if}
      </button>
    {/each}
  </div>
  <span class="hint" class:custom={active === null}>{hint}</span>
  <span class="count">{S.profiles.selected(selectedCount, fixes.length)}</span>
</div>

<style>
  .bar {
    position: sticky;
    top: 0;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 10px 2px;
    background: color-mix(in srgb, var(--bg-window) 90%, transparent);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid var(--stroke);
  }
  .label {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-faint);
  }
  .segments {
    display: flex;
    gap: 2px;
    padding: 3px;
    background: rgba(0, 0, 0, 0.22);
    border: 1px solid var(--stroke);
    border-radius: var(--r-control);
    box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.28);
  }
  .segment {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 5px 12px;
    border-radius: 4px;
    font-size: 12.5px;
    color: var(--text-dim);
    white-space: nowrap;
    transition:
      background 0.12s ease,
      color 0.12s ease;
  }
  .segment:hover {
    color: var(--text);
    background: var(--bg-raised);
  }
  .segment.active {
    background: linear-gradient(180deg, var(--coral-bright), var(--coral));
    color: #241511;
    font-weight: 600;
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.22),
      0 1px 2px rgba(0, 0, 0, 0.3);
  }
  .dot {
    width: 5px;
    height: 5px;
    border-radius: 50%;
    background: var(--gold);
  }
  .segment.active .dot {
    background: #241511;
  }
  .hint {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 12px;
    color: var(--text-faint);
  }
  .hint.custom {
    color: var(--gold);
  }
  .count {
    flex: none;
    font-size: 12px;
    color: var(--text-dim);
    font-variant-numeric: tabular-nums;
  }
</style>
