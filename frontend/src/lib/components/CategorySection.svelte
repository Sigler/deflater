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

  type Item =
    | { kind: "single"; fix: FixState }
    | { kind: "group"; id: string; primary: FixState; children: FixState[] };

  // Cluster grouped fixes into one card at the position of the group's
  // first member; everything else renders as a standalone row.
  const items = $derived.by<Item[]>(() => {
    const out: Item[] = [];
    const seen = new Set<string>();
    for (const fix of fixes) {
      if (!fix.group) {
        out.push({ kind: "single", fix });
        continue;
      }
      if (seen.has(fix.group)) continue;
      seen.add(fix.group);
      const members = fixes.filter((f) => f.group === fix.group);
      const primary = members.find((m) => m.primary) ?? members[0];
      out.push({
        kind: "group",
        id: fix.group,
        primary,
        children: members.filter((m) => m !== primary),
      });
    }
    return out;
  });
</script>

{#if fixes.length > 0}
  <section id={"sec-" + id}>
    <header>
      <h2>{text?.title ?? id}</h2>
      <p>{text?.blurb ?? ""}</p>
      <span class="count">{S.profiles.selected(onCount, fixes.length)}</span>
    </header>
    <div class="list">
      {#each items as item (item.kind === "single" ? item.fix.id : item.id)}
        {#if item.kind === "single"}
          <FixRow
            fix={item.fix}
            selected={selection.has(item.fix.id)}
            pending={pending.has(item.fix.id)}
            {applying}
            {ontoggle}
          />
        {:else}
          <div class="group">
            <FixRow
              fix={item.primary}
              selected={selection.has(item.primary.id)}
              pending={pending.has(item.primary.id)}
              {applying}
              {ontoggle}
              flat
            />
            {#each item.children as childFix (childFix.id)}
              <div class="childsep"></div>
              <FixRow
                fix={childFix}
                selected={selection.has(childFix.id)}
                pending={pending.has(childFix.id)}
                {applying}
                {ontoggle}
                flat
                child
              />
            {/each}
          </div>
        {/if}
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
  /* One card wrapping a primary fix and its sub-options. */
  .group {
    background: var(--bg-panel);
    border: 1px solid var(--stroke);
    border-radius: var(--r-card);
    overflow: hidden;
    transition: border-color 0.12s ease;
  }
  .group:hover {
    border-color: var(--stroke-strong);
  }
  .childsep {
    height: 1px;
    background: var(--stroke);
  }
</style>
