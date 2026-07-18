<script lang="ts">
  import { S } from "../i18n";
  import type { FixState } from "../types";
  import { matchesProfile } from "../changes";

  let {
    fixes,
    selection,
    onpick,
  }: {
    fixes: FixState[];
    selection: Set<string>;
    onpick: (profile: string) => void;
  } = $props();

  const profiles = ["light-touch", "clean-sweep", "full-deflate"] as const;

  const active = $derived(profiles.find((p) => matchesProfile(fixes, selection, p)) ?? null);

  function count(profile: string): number {
    return fixes.filter((f) => f.profiles.includes(profile)).length;
  }
</script>

<div class="wrap">
  <div class="heading">
    <h2>{S.profiles.heading}</h2>
    <p>
      {S.profiles.subheading}
      {#if active === null}<span class="custom">{S.profiles.custom}</span>{/if}
    </p>
  </div>
  <div class="cards">
    {#each profiles as p (p)}
      <button type="button" class="card" class:active={active === p} onclick={() => onpick(p)}>
        <span class="top">
          <span class="name">{S.profiles[p].title}</span>
          {#if p === "clean-sweep"}
            <span class="badge">{S.profiles["clean-sweep"].badge}</span>
          {/if}
        </span>
        <span class="tagline">{S.profiles[p].tagline}</span>
        <span class="count">{count(p)} fixes</span>
      </button>
    {/each}
  </div>
</div>

<style>
  .wrap {
    display: grid;
    gap: 12px;
  }
  .heading {
    display: flex;
    align-items: baseline;
    gap: 12px;
    padding: 0 2px;
  }
  h2 {
    font-size: 16px;
    font-weight: 600;
  }
  .heading p {
    color: var(--text-faint);
    font-size: 12.5px;
  }
  .custom {
    margin-left: 8px;
    color: var(--gold);
  }
  .cards {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }
  .card {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    text-align: left;
    padding: 14px 16px;
    background: var(--bg-card);
    border: 1px solid var(--stroke);
    border-radius: var(--r-card);
    transition:
      border-color 0.12s ease,
      background 0.12s ease;
  }
  .card:hover {
    border-color: var(--stroke-strong);
    background: var(--bg-raised);
  }
  .card.active {
    border-color: var(--coral);
    background: var(--coral-soft);
  }
  .top {
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .name {
    font-weight: 600;
    font-size: 14px;
  }
  .badge {
    font-size: 10.5px;
    line-height: 1;
    padding: 4px 8px;
    border-radius: var(--r-chip);
    background: var(--gold-soft);
    color: var(--gold);
  }
  .tagline {
    font-size: 12.5px;
    line-height: 1.45;
    color: var(--text-dim);
  }
  .count {
    font-size: 11.5px;
    color: var(--text-faint);
  }
</style>
