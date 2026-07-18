<script lang="ts">
  import { fly } from "svelte/transition";
  import { S } from "../i18n";

  let {
    changeCount,
    applying,
    progressText,
    onapply,
    onreset,
  }: {
    changeCount: number;
    applying: boolean;
    progressText: string;
    onapply: () => void;
    onreset: () => void;
  } = $props();
</script>

{#if changeCount > 0 || applying}
  <div class="bar" transition:fly={{ y: 28, duration: 220 }}>
    <span class="count">
      {applying ? progressText || S.apply.applying : S.apply.changesReady(changeCount)}
    </span>
    <div class="actions">
      <button type="button" class="ghost" disabled={applying} onclick={onreset}>
        {S.apply.reset}
      </button>
      <button type="button" class="primary" disabled={applying} onclick={onapply}>
        {applying ? S.apply.applying : S.apply.applyCount(changeCount)}
      </button>
    </div>
  </div>
{/if}

<style>
  .bar {
    position: sticky;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 12px 20px;
    background: color-mix(in srgb, var(--bg-panel) 88%, transparent);
    backdrop-filter: blur(12px);
    border-top: 1px solid var(--stroke);
  }
  .count {
    font-size: 13px;
    color: var(--text-dim);
  }
  .actions {
    display: flex;
    gap: 10px;
  }
  .ghost {
    padding: 8px 16px;
    border-radius: var(--r-control);
    border: 1px solid var(--stroke-strong);
    color: var(--text-dim);
  }
  .ghost:hover:not(:disabled) {
    color: var(--text);
    border-color: var(--text-faint);
  }
  .primary {
    padding: 8px 22px;
    border-radius: var(--r-control);
    background: linear-gradient(180deg, var(--coral-bright), var(--coral));
    color: #241511;
    font-weight: 600;
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.25),
      0 1px 3px rgba(0, 0, 0, 0.35);
    transition:
      background 0.12s ease,
      box-shadow 0.12s ease,
      transform 0.06s ease;
  }
  .primary:hover:not(:disabled) {
    background: linear-gradient(180deg, #ff8f70, var(--coral-bright));
  }
  .primary:active:not(:disabled) {
    transform: translateY(1px);
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.15),
      0 1px 1px rgba(0, 0, 0, 0.3);
  }
  .primary:disabled {
    opacity: 0.6;
  }
</style>
