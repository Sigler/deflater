<script lang="ts">
  import mascot from "../../assets/mascot-512.png";
  import { S } from "../i18n";
  import LogoMark from "./LogoMark.svelte";

  let {
    elevated,
    inPlace,
    total,
  }: { elevated: boolean; inPlace: number; total: number } = $props();
</script>

<header>
  <div class="brand">
    <LogoMark size={44} />
    <div class="words">
      <h1>{S.app.name}</h1>
      <p class="tagline">{S.app.tagline}</p>
      <p class="readout">{S.app.readout(inPlace, total)}</p>
      <p class="assurance">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path
            d="M12 2 L20 5.5 V11 C20 16.5 16.6 20.6 12 22 C7.4 20.6 4 16.5 4 11 V5.5 Z"
            stroke="var(--sage)"
            stroke-width="2"
            stroke-linejoin="round"
          />
          <path
            d="M8.5 11.5 L11 14 L15.5 9"
            stroke="var(--sage)"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        {S.app.assurance}
      </p>
    </div>
  </div>
  <div class="side">
    <span class="mode" class:admin={elevated} title={elevated ? "" : S.app.notElevatedHint}>
      {elevated ? S.app.elevated : S.app.notElevated}
    </span>
    <img class="mascot" src={mascot} alt="" draggable="false" />
  </div>
</header>

<style>
  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 20px;
  }
  .brand {
    display: flex;
    align-items: center;
    gap: 14px;
    min-width: 0;
  }
  .words {
    display: grid;
    gap: 3px;
    min-width: 0;
  }
  h1 {
    font-size: 22px;
    font-weight: 700;
    letter-spacing: -0.01em;
  }
  .tagline {
    color: var(--text-dim);
    font-size: 13px;
  }
  .readout {
    color: var(--text-faint);
    font-size: 12px;
    font-variant-numeric: tabular-nums;
  }
  .assurance {
    display: flex;
    align-items: center;
    gap: 7px;
    margin-top: 4px;
    font-size: 12px;
    line-height: 1.5;
    color: var(--text-dim);
  }
  .assurance svg {
    flex: none;
  }
  .side {
    flex: none;
    display: flex;
    align-items: center;
    gap: 14px;
  }
  .mode {
    font-size: 11.5px;
    padding: 5px 10px;
    border-radius: var(--r-chip);
    background: var(--bg-raised);
    color: var(--text-dim);
    border: 1px solid var(--stroke);
    white-space: nowrap;
  }
  .mode.admin {
    color: var(--sage);
    background: var(--sage-soft);
  }
  .mascot {
    width: 132px;
    height: 132px;
    border-radius: 10px;
    border: 1px solid var(--stroke-strong);
    box-shadow: 0 6px 18px rgba(0, 0, 0, 0.45);
    transform: rotate(-1.2deg);
    user-select: none;
  }
</style>
