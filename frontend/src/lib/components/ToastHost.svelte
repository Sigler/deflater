<script lang="ts">
  import { fly } from "svelte/transition";
  import { S } from "../i18n";
  import { dismissToast, pushToast, type Toast, toasts } from "../toasts.svelte";

  let busy = $state<number | null>(null);

  async function runAction(t: Toast) {
    if (!t.action) return;
    busy = t.id;
    try {
      await t.action.run();
      dismissToast(t.id);
    } catch (e) {
      // Don't dismiss as if it worked; tell the user it failed.
      pushToast({ kind: "warn", message: `${e}`, sticky: true });
    } finally {
      busy = null;
    }
  }
</script>

<div class="host" aria-live="polite" aria-atomic="false">
  {#each toasts as t (t.id)}
    <div class="toast {t.kind}" role="status" transition:fly={{ y: 16, duration: 200 }}>
      <span class="dot" aria-hidden="true"></span>
      <div class="body">
        <p class="msg">{t.message}</p>
        {#if t.detail}
          {#each t.detail as line (line)}
            <p class="detail">{line}</p>
          {/each}
        {/if}
        {#if t.action}
          <button type="button" class="action" disabled={busy === t.id} onclick={() => runAction(t)}>
            {busy === t.id ? (t.action.busyLabel ?? t.action.label) : t.action.label}
          </button>
        {/if}
      </div>
      <button type="button" class="x" onclick={() => dismissToast(t.id)} aria-label={S.toast.dismiss}>
        ×
      </button>
    </div>
  {/each}
</div>

<style>
  .host {
    position: fixed;
    right: 20px;
    bottom: 76px;
    z-index: 50;
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: min(360px, calc(100vw - 40px));
    pointer-events: none;
  }
  .toast {
    pointer-events: auto;
    display: grid;
    grid-template-columns: auto 1fr auto;
    align-items: start;
    gap: 10px;
    padding: 11px 12px 11px 13px;
    background: var(--bg-raised);
    border: 1px solid var(--stroke-strong);
    border-radius: var(--r-card);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  }
  .dot {
    margin-top: 5px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--text-faint);
  }
  .success .dot {
    background: var(--sage);
  }
  .warn .dot {
    background: var(--gold);
  }
  .info .dot {
    background: var(--coral);
  }
  .body {
    min-width: 0;
    display: grid;
    gap: 3px;
  }
  .msg {
    font-size: 13px;
    color: var(--text);
  }
  .detail {
    font-size: 12px;
    color: var(--text-dim);
  }
  .action {
    justify-self: start;
    margin-top: 4px;
    font-size: 12px;
    padding: 5px 12px;
    border-radius: var(--r-control);
    border: 1px solid var(--stroke-strong);
    background: var(--bg-card);
    color: var(--text);
  }
  .action:hover:not(:disabled) {
    border-color: var(--coral);
    color: var(--coral-bright);
  }
  .action:disabled {
    opacity: 0.6;
  }
  .x {
    color: var(--text-faint);
    font-size: 16px;
    line-height: 1;
    padding: 2px 4px;
  }
  .x:hover {
    color: var(--text);
  }
</style>
