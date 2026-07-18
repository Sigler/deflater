<script lang="ts">
  import { S } from "../i18n";
  import type { Alert } from "../types";

  let {
    alerts,
    onremove,
    ondismiss,
  }: {
    alerts: Alert[];
    onremove: (pkg: string) => Promise<void>;
    ondismiss: () => void;
  } = $props();

  let removing = $state<string | null>(null);
  let error = $state("");

  async function remove(pkg: string) {
    removing = pkg;
    error = "";
    try {
      await onremove(pkg);
    } catch (e) {
      error = `${pkg}: ${e}`;
    } finally {
      removing = null;
    }
  }
</script>

{#if alerts.length > 0}
  <div class="banner" role="alert">
    <div class="head">
      <span class="title">{S.alerts.title}</span>
      <button type="button" class="dismiss" onclick={ondismiss}>{S.alerts.dismiss}</button>
    </div>
    <p>{S.alerts.body}</p>
    <ul>
      {#each alerts as alert (alert.package)}
        <li>
          <code>{alert.package}</code>
          <button
            type="button"
            class="remove"
            disabled={removing !== null}
            aria-label={`${S.alerts.remove} ${alert.package}`}
            onclick={() => remove(alert.package)}
          >
            {removing === alert.package ? S.alerts.removing : S.alerts.remove}
          </button>
        </li>
      {/each}
    </ul>
    {#if error}<p class="err">{error}</p>{/if}
  </div>
{/if}

<style>
  .banner {
    background: var(--gold-soft);
    border: 1px solid rgba(214, 164, 76, 0.35);
    border-radius: var(--r-card);
    padding: 14px 16px;
    display: grid;
    gap: 8px;
  }
  .head {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .title {
    font-weight: 600;
    color: var(--gold);
  }
  .dismiss {
    font-size: 12px;
    color: var(--text-dim);
  }
  .dismiss:hover {
    color: var(--text);
  }
  p {
    font-size: 12.5px;
    color: var(--text-dim);
  }
  .err {
    color: var(--coral-bright);
  }
  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 6px;
  }
  li {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
  }
  code {
    font-size: 12.5px;
    color: var(--text);
  }
  .remove {
    font-size: 12px;
    padding: 5px 12px;
    border-radius: var(--r-control);
    background: var(--bg-raised);
    border: 1px solid var(--stroke-strong);
  }
  .remove:hover:not(:disabled) {
    border-color: var(--coral);
    color: var(--coral-bright);
  }
</style>
