<script lang="ts">
  import { slide } from "svelte/transition";
  import { S } from "../i18n";
  import type { FixState } from "../types";
  import StatusChip from "./StatusChip.svelte";
  import Toggle from "./Toggle.svelte";

  let {
    fix,
    selected,
    pending,
    ontoggle,
  }: {
    fix: FixState;
    selected: boolean;
    pending: boolean;
    ontoggle: (id: string, value: boolean) => void;
  } = $props();

  let expanded = $state(false);

  const text = $derived(S.fixes[fix.id as keyof typeof S.fixes]);
  const isApp = $derived(fix.kind === "app-junk" || fix.kind === "app-might");
  // A gone app cannot be brought back by a toggle; lock it on.
  const locked = $derived(isApp && fix.status === "removed");

  const mechanism = $derived.by(() => {
    const parts: string[] = [];
    if (fix.appx?.length) parts.push(S.details.mechanismApp(fix.appx.join(", ")));
    if (fix.reg?.length) parts.push(S.details.mechanismReg(fix.reg.length));
    return parts.join(" ");
  });
</script>

<div class="row" class:expanded>
  <div class="main">
    <button
      type="button"
      class="info"
      aria-expanded={expanded}
      onclick={() => (expanded = !expanded)}
    >
      <span class="titleline">
        <span class="title">{text?.title ?? fix.id}</span>
        <svg
          class="chevron"
          class:open={expanded}
          width="10"
          height="10"
          viewBox="0 0 10 10"
          aria-hidden="true"
        >
          <path
            d="M2 3.5 L5 6.5 L8 3.5"
            fill="none"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        {#if fix.caution}
          <span class="caution">{S.badges.caution}</span>
        {/if}
      </span>
      <span class="summary">{text?.summary ?? ""}</span>
    </button>
    <StatusChip status={fix.status} kind={fix.kind} />
    {#if pending}
      <span class="pending" title={S.badges.willChange} aria-label={S.badges.willChange}></span>
    {/if}
    <Toggle
      checked={selected}
      disabled={locked}
      label={text?.title ?? fix.id}
      onchange={(v) => ontoggle(fix.id, v)}
    />
  </div>

  {#if expanded && text}
    <div class="details" transition:slide={{ duration: 180 }}>
      <div class="block">
        <span class="label">{S.details.what}</span>
        <p>{text.what}</p>
        {#if mechanism}<p class="mech">{mechanism}</p>{/if}
      </div>
      {#if text.tradeoff}
        <div class="block">
          <span class="label gold">{S.details.tradeoff}</span>
          <p>{text.tradeoff}</p>
        </div>
      {/if}
      <div class="block">
        <span class="label">{S.details.undo}</span>
        <p>{isApp ? S.details.undoApp : text.undo}</p>
      </div>
    </div>
  {/if}
</div>

<style>
  .row {
    background:
      linear-gradient(180deg, rgba(255, 244, 230, 0.025), transparent 45%),
      var(--bg-card);
    border: 1px solid var(--stroke);
    border-radius: var(--r-card);
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.22);
    transition:
      border-color 0.12s ease,
      box-shadow 0.12s ease;
  }
  .row:hover {
    border-color: var(--stroke-strong);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.28);
  }
  .main {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 12px 16px;
  }
  .info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 3px;
    text-align: left;
    min-width: 0;
  }
  .titleline {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .title {
    font-weight: 600;
    font-size: 14px;
  }
  .chevron {
    flex: none;
    color: var(--text-faint);
    transition:
      transform 0.15s ease,
      color 0.12s ease;
  }
  .row:hover .chevron {
    color: var(--text-dim);
  }
  .chevron.open {
    transform: rotate(180deg);
    color: var(--coral);
  }
  .caution {
    font-size: 11px;
    line-height: 1;
    padding: 4px 8px;
    border-radius: var(--r-chip);
    background: var(--gold-soft);
    color: var(--gold);
  }
  .summary {
    color: var(--text-dim);
    font-size: 12.5px;
  }
  .pending {
    flex: none;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--coral);
    box-shadow: 0 0 6px rgba(229, 106, 77, 0.55);
  }
  .details {
    border-top: 1px solid var(--stroke);
    padding: 14px 16px 16px;
    display: grid;
    gap: 12px;
  }
  .block {
    display: grid;
    gap: 4px;
  }
  .label {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--text-faint);
  }
  .label.gold {
    color: var(--gold);
  }
  .details p {
    font-size: 13px;
    line-height: 1.55;
    color: var(--text-dim);
    max-width: 68ch;
  }
  .mech {
    font-size: 12px !important;
    color: var(--text-faint) !important;
  }
</style>
