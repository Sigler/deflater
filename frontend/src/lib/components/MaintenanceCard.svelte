<script lang="ts">
  import { S } from "../i18n";
  import Toggle from "./Toggle.svelte";

  let {
    maintenance,
    maintenancePending,
    onmaintenance,
  }: {
    maintenance: boolean;
    maintenancePending: boolean;
    onmaintenance: (on: boolean) => void;
  } = $props();
</script>

<div class="card">
  <div class="line">
    <div class="text">
      <span class="title">{S.maintenance.title}</span>
      <span class="body">{S.maintenance.body}</span>
      {#if maintenance && maintenancePending}
        <span class="pending">{S.maintenance.pendingElevation}</span>
      {/if}
    </div>
    <Toggle checked={maintenance} label={S.maintenance.title} onchange={onmaintenance} />
  </div>
</div>

<style>
  .card {
    background: var(--bg-panel);
    border: 1px solid var(--stroke);
    border-radius: var(--r-card);
    padding: 4px 16px;
  }
  .line {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 0;
  }
  .text {
    flex: 1;
    display: grid;
    gap: 3px;
  }
  .title {
    font-weight: 600;
    font-size: 14px;
  }
  .body {
    font-size: 12.5px;
    line-height: 1.5;
    color: var(--text-dim);
    max-width: 72ch;
  }
  .pending {
    font-size: 12px;
    color: var(--gold);
  }
</style>
