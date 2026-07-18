<script lang="ts">
  import { S } from "../i18n";
  import type { RecallInfo } from "../types";

  let { info, onopen }: { info: RecallInfo; onopen: () => void } = $props();

  function formatBytes(n: number): string {
    if (n < 1024) return `${n} B`;
    const units = ["KB", "MB", "GB", "TB"];
    let v = n / 1024;
    let i = 0;
    while (v >= 1024 && i < units.length - 1) {
      v /= 1024;
      i++;
    }
    return `${v >= 10 ? Math.round(v) : v.toFixed(1)} ${units[i]}`;
  }
</script>

{#if info.present}
  <div class="banner">
    <span class="dot" aria-hidden="true"></span>
    <div class="body">
      <span class="title">{S.recall.title}</span>
      <p>{S.recall.body(formatBytes(info.bytes))}</p>
      <div class="row">
        <code>{info.path}</code>
        <button type="button" class="open" onclick={onopen}>{S.recall.openFolder}</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .banner {
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 10px;
    padding: 13px 15px;
    background: var(--bg-card);
    border: 1px solid var(--stroke-strong);
    border-radius: var(--r-card);
  }
  .dot {
    margin-top: 6px;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--coral);
  }
  .body {
    min-width: 0;
    display: grid;
    gap: 6px;
  }
  .title {
    font-weight: 600;
    color: var(--text);
  }
  p {
    font-size: 12.5px;
    color: var(--text-dim);
  }
  .row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    flex-wrap: wrap;
  }
  code {
    font-size: 11.5px;
    color: var(--text-faint);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }
  .open {
    flex: none;
    font-size: 12px;
    padding: 5px 12px;
    border-radius: var(--r-control);
    background: var(--bg-raised);
    border: 1px solid var(--stroke-strong);
    color: var(--text);
  }
  .open:hover {
    border-color: var(--coral);
    color: var(--coral-bright);
  }
</style>
