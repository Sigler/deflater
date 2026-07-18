<script lang="ts">
  import { S } from "../i18n";

  let { inPlace, total }: { inPlace: number; total: number } = $props();
</script>

<div class="bar">
  <span class="what">
    <svg width="13" height="13" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path
        d="M4 12.5 L9.5 18 L20 6.5"
        stroke="var(--sage)"
        stroke-width="2.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
    {S.app.scanChecked}
  </span>
  <span class="count">{S.app.scanApplied(inPlace, total)}</span>
</div>

<style>
  /* Stickiness and backdrop belong to the wrapper in App.svelte, so the
     tabs can share the same sticky block. */
  .bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    padding: 10px 2px 6px;
  }
  .what {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 12.5px;
    color: var(--text-dim);
  }
  .what svg {
    flex: none;
  }
  .count {
    flex: none;
    font-size: 12.5px;
    color: var(--text);
    font-variant-numeric: tabular-nums;
  }
</style>
