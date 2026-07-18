<script lang="ts">
  import { S } from "../i18n";
  import Toggle from "./Toggle.svelte";

  let {
    maintenance,
    watcher,
    pendingElevation,
    onmaintenance,
    onwatcher,
  }: {
    maintenance: boolean;
    watcher: boolean;
    pendingElevation: boolean;
    onmaintenance: (on: boolean) => void;
    onwatcher: (on: boolean) => void;
  } = $props();
</script>

<div class="card">
  <div class="line">
    <div class="text">
      <span class="title">{S.maintenance.title}</span>
      <span class="body">{S.maintenance.body}</span>
      {#if maintenance && pendingElevation}
        <span class="pending">{S.maintenance.pendingElevation}</span>
      {/if}
    </div>
    <Toggle checked={maintenance} label={S.maintenance.title} onchange={onmaintenance} />
  </div>
  <div class="line sub" class:dimmed={!maintenance}>
    <div class="text">
      <span class="title small">{S.maintenance.watcherTitle}</span>
      <span class="body">{S.maintenance.watcherBody}</span>
    </div>
    <Toggle
      checked={watcher}
      disabled={!maintenance}
      label={S.maintenance.watcherTitle}
      onchange={onwatcher}
    />
  </div>
</div>

<style>
  .card {
    background:
      linear-gradient(180deg, rgba(214, 228, 255, 0.025), transparent 45%),
      var(--bg-panel);
    border: 1px solid var(--stroke);
    border-radius: var(--r-card);
    padding: 4px 16px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.22);
  }
  .line {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 12px 0;
  }
  .line.sub {
    border-top: 1px solid var(--stroke);
  }
  .line.sub.dimmed .text {
    opacity: 0.5;
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
  .title.small {
    font-size: 13px;
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
