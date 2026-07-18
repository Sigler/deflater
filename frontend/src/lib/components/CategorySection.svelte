<script lang="ts">
  import { S } from "../i18n";
  import type { FixState } from "../types";
  import FixRow from "./FixRow.svelte";

  let {
    id,
    fixes,
    selection,
    pending,
    applying,
    ontoggle,
  }: {
    id: string;
    fixes: FixState[];
    selection: Set<string>;
    pending: Set<string>;
    applying: boolean;
    ontoggle: (id: string, value: boolean) => void;
  } = $props();

  const text = $derived(S.categories[id as keyof typeof S.categories]);
  const onCount = $derived(fixes.filter((f) => selection.has(f.id)).length);
</script>

{#if fixes.length > 0}
  <section id={"sec-" + id}>
    <header>
      <h2>{text?.title ?? id}</h2>
      <p>{text?.blurb ?? ""}</p>
      <span class="count">{S.profiles.selected(onCount, fixes.length)}</span>
    </header>
    <div class="list">
      {#each fixes as fix (fix.id)}
        <FixRow
          {fix}
          selected={selection.has(fix.id)}
          pending={pending.has(fix.id)}
          {applying}
          {ontoggle}
        />
      {/each}
    </div>
  </section>
{/if}

<style>
  section {
    display: grid;
    gap: 10px;
    scroll-margin-top: 96px;
    /* Extra air between crates so sections read as distinct blocks.
       Adds to the parent grid gap; the tight 8-10px internal gaps stay. */
    margin-top: 18px;
  }
  /* The first section sits right under the sticky bar; no double gap. */
  section:first-of-type {
    margin-top: 4px;
  }
  header {
    display: flex;
    align-items: baseline;
    gap: 12px;
    padding: 0 2px;
  }
  /* Stencil headers: uppercase, spaced, like markings on a crate. */
  h2 {
    font-size: 13px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.09em;
    color: var(--text);
  }
  header p {
    color: var(--text-faint);
    font-size: 12.5px;
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .count {
    flex: none;
    font-size: 11.5px;
    color: var(--text-faint);
    font-variant-numeric: tabular-nums;
  }
  .list {
    display: grid;
    gap: 8px;
  }
</style>
