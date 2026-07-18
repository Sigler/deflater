<script lang="ts">
  import { S } from "../i18n";
  import type { FixState } from "../types";
  import FixRow from "./FixRow.svelte";

  let {
    id,
    fixes,
    selection,
    ontoggle,
  }: {
    id: string;
    fixes: FixState[];
    selection: Set<string>;
    ontoggle: (id: string, value: boolean) => void;
  } = $props();

  const text = $derived(S.categories[id as keyof typeof S.categories]);
</script>

{#if fixes.length > 0}
  <section>
    <header>
      <h2>{text?.title ?? id}</h2>
      <p>{text?.blurb ?? ""}</p>
    </header>
    <div class="list">
      {#each fixes as fix (fix.id)}
        <FixRow {fix} selected={selection.has(fix.id)} {ontoggle} />
      {/each}
    </div>
  </section>
{/if}

<style>
  section {
    display: grid;
    gap: 10px;
  }
  header {
    display: flex;
    align-items: baseline;
    gap: 12px;
    padding: 0 2px;
  }
  h2 {
    font-size: 16px;
    font-weight: 600;
  }
  header p {
    color: var(--text-faint);
    font-size: 12.5px;
  }
  .list {
    display: grid;
    gap: 8px;
  }
</style>
