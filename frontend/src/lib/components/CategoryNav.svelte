<script lang="ts">
  export interface NavItem {
    id: string;
    label: string;
  }

  let {
    items,
    active,
    variant,
    onjump,
  }: {
    items: NavItem[];
    active: string;
    variant: "tabs" | "side";
    onjump: (id: string) => void;
  } = $props();
</script>

<nav class={variant} aria-label="Sections">
  {#each items as item (item.id)}
    <button
      type="button"
      class:active={active === item.id}
      aria-current={active === item.id ? "true" : undefined}
      onclick={() => onjump(item.id)}
    >
      {item.label}
    </button>
  {/each}
</nav>

<style>
  nav.tabs {
    display: flex;
    gap: 2px;
    overflow-x: auto;
    scrollbar-width: none;
  }
  nav.tabs::-webkit-scrollbar {
    display: none;
  }
  nav.tabs button {
    flex: none;
    font-size: 12px;
    padding: 7px 10px 9px;
    color: var(--text-dim);
    border-bottom: 2px solid transparent;
    white-space: nowrap;
    transition:
      color 0.12s ease,
      border-color 0.12s ease;
  }
  nav.tabs button:hover {
    color: var(--text);
  }
  nav.tabs button.active {
    color: var(--coral-bright);
    border-bottom-color: var(--coral);
  }

  nav.side {
    display: grid;
    gap: 2px;
  }
  nav.side button {
    text-align: left;
    font-size: 12.5px;
    padding: 7px 10px;
    border-radius: var(--r-control);
    color: var(--text-dim);
    transition:
      color 0.12s ease,
      background 0.12s ease;
  }
  nav.side button:hover {
    color: var(--text);
    background: var(--bg-panel);
  }
  nav.side button.active {
    color: var(--coral-bright);
    background: var(--coral-soft);
  }
</style>
