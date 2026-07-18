<script lang="ts">
  import { S } from "../i18n";
  import type { ForeignTask } from "../types";

  let {
    tasks,
    onremove,
  }: {
    tasks: ForeignTask[];
    // Resolves once the removal is done (or has been handed off to an
    // elevated relaunch, in which case this window is about to close).
    onremove: (name: string) => Promise<void>;
  } = $props();

  let removing = $state<string | null>(null);
  let error = $state("");

  async function remove(name: string) {
    removing = name;
    error = "";
    try {
      await onremove(name);
    } catch (e) {
      error = `${e}`;
    } finally {
      removing = null;
    }
  }
</script>

{#if tasks.length > 0}
  <div class="banner" role="alert">
    <span class="title">{S.conflicts.title}</span>
    <p>{S.conflicts.body}</p>
    <ul>
      {#each tasks as task (task.name)}
        <li>
          <div class="what">
            <code>{task.name}</code>
            <span class="tool">{S.conflicts.fromTool(task.tool)}</span>
            {#if task.note}<span class="note">{task.note}</span>{/if}
          </div>
          <button
            type="button"
            class="remove"
            disabled={removing !== null}
            aria-label={`${S.conflicts.remove}: ${task.name}`}
            onclick={() => remove(task.name)}
          >
            {removing === task.name ? S.conflicts.removing : S.conflicts.remove}
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
  .title {
    font-weight: 600;
    color: var(--gold);
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
  .what {
    display: grid;
    gap: 2px;
    min-width: 0;
  }
  code {
    font-size: 12.5px;
    color: var(--text);
  }
  .tool {
    font-size: 11.5px;
    color: var(--text-faint);
  }
  .note {
    font-size: 11.5px;
    color: var(--text-dim);
  }
  .remove {
    flex: none;
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
