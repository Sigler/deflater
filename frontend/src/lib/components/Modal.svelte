<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    title,
    children,
    actions,
  }: {
    title: string;
    children: Snippet;
    actions: Snippet;
  } = $props();
</script>

<div class="overlay">
  <div class="modal" role="dialog" aria-modal="true" aria-label={title}>
    <h3>{title}</h3>
    <div class="body">{@render children()}</div>
    <div class="actions">{@render actions()}</div>
  </div>
</div>

<style>
  .overlay {
    position: fixed;
    inset: 0;
    background: rgba(7, 10, 16, 0.62);
    display: grid;
    place-items: center;
    z-index: 50;
  }
  .modal {
    width: min(440px, calc(100vw - 48px));
    background: var(--bg-panel);
    border: 1px solid var(--stroke-strong);
    border-radius: 10px;
    padding: 20px;
    display: grid;
    gap: 14px;
    box-shadow: 0 18px 50px rgba(0, 0, 0, 0.5);
  }
  h3 {
    font-size: 16px;
    font-weight: 600;
  }
  .body {
    font-size: 13px;
    line-height: 1.55;
    color: var(--text-dim);
  }
  .actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }
</style>
