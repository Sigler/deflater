<script lang="ts">
  import mascot from "../../assets/mascot-512.png";
  import { S } from "../i18n";

  // The mascot is small in the masthead; clicking it is the easter egg.
  let showBig = $state(false);
</script>

<svelte:window
  onkeydown={(e) => {
    if (e.key === "Escape") showBig = false;
  }}
/>

<header>
  <button
    type="button"
    class="mascot-btn"
    aria-label="Deflater"
    onclick={() => (showBig = true)}
  >
    <img class="mascot" src={mascot} alt="" draggable="false" />
  </button>
  <div class="words">
    <h1>{S.app.name}</h1>
    <p class="tagline">{S.app.tagline}</p>
  </div>
</header>

{#if showBig}
  <button type="button" class="lightbox" aria-label="Close" onclick={() => (showBig = false)}>
    <img src={mascot} alt="The Deflater mascot" draggable="false" />
  </button>
{/if}

<style>
  header {
    display: flex;
    align-items: center;
    gap: 18px;
  }
  .mascot-btn {
    flex: none;
    line-height: 0;
    cursor: pointer;
  }
  .mascot {
    width: 74px;
    height: 74px;
    user-select: none;
    transition: transform 0.15s ease;
  }
  .mascot-btn:hover .mascot {
    transform: scale(1.06) rotate(-1.5deg);
  }
  .words {
    display: grid;
    gap: 2px;
    min-width: 0;
  }
  h1 {
    font-size: 24px;
    font-weight: 700;
    letter-spacing: -0.01em;
    /* Misregistered print: a coral pass slightly off the black one. */
    text-shadow: 2px 2px 0 rgba(229, 106, 77, 0.35);
  }
  .tagline {
    color: var(--text-dim);
    font-size: 13px;
  }
  .lightbox {
    position: fixed;
    inset: 0;
    z-index: 60;
    display: grid;
    place-items: center;
    background: rgba(7, 10, 16, 0.82);
    cursor: zoom-out;
  }
  .lightbox img {
    width: min(78vw, 68vh, 560px);
    height: auto;
    filter: drop-shadow(0 24px 60px rgba(0, 0, 0, 0.6));
  }
</style>
