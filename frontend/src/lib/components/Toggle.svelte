<script lang="ts">
  let {
    checked,
    disabled = false,
    label,
    onchange,
  }: {
    checked: boolean;
    disabled?: boolean;
    label: string;
    onchange: (value: boolean) => void;
  } = $props();
</script>

<button
  type="button"
  role="switch"
  aria-checked={checked}
  aria-label={label}
  class="toggle"
  class:on={checked}
  {disabled}
  onclick={() => onchange(!checked)}
>
  <span class="knob"></span>
</button>

<style>
  /* Soft-UI toggle, one light source at the upper left.
     The track is a well sunk into the card: its upper-left inner lip is
     in shadow, its lower-right inner edge catches faint reflected light,
     and a soft rim highlight sits on the outer lower edge where the
     surface rolls back up. The state color lies matte at the bottom of
     the well. The knob is a soft dome resting proud of the well: bright
     toward the light at the upper left, curving into shade at the lower
     right, dropping a diffuse shadow onto the well floor. */
  .toggle {
    flex: none;
    width: 48px;
    height: 26px;
    border-radius: 999px;
    border: none;
    background: linear-gradient(150deg, rgba(0, 0, 0, 0.42), rgba(0, 0, 0, 0.18));
    box-shadow:
      inset 0 0 0 1px rgba(0, 0, 0, 0.55),
      inset 2px 3px 5px rgba(0, 0, 0, 0.48),
      inset -1px -2px 3px rgba(214, 228, 255, 0.05),
      1px 2px 3px rgba(214, 228, 255, 0.05);
    position: relative;
    transition: background 0.18s ease;
  }
  .toggle.on {
    background: linear-gradient(150deg, #8a3626, #c05539 55%, #cd5f40);
    box-shadow:
      inset 0 0 0 1px rgba(0, 0, 0, 0.38),
      inset 2px 3px 5px rgba(0, 0, 0, 0.42),
      inset -1px -2px 3px rgba(255, 200, 170, 0.1),
      1px 2px 3px rgba(214, 228, 255, 0.05);
  }
  .toggle.on:hover:not(:disabled) {
    background: linear-gradient(150deg, #964130, #cd5f40 55%, #d96a49);
  }
  .toggle:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .knob {
    position: absolute;
    top: 3px;
    left: 3px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: radial-gradient(circle at 32% 28%, #f3f6fa, #d3dae4 55%, #aeb8c6);
    box-shadow:
      inset 1px 1px 1px rgba(255, 255, 255, 0.55),
      inset -2px -3px 4px rgba(0, 0, 0, 0.16),
      2px 3px 5px rgba(0, 0, 0, 0.4),
      1px 1px 2px rgba(0, 0, 0, 0.35);
    transition:
      transform 0.18s cubic-bezier(0.3, 0.9, 0.4, 1.1),
      background 0.15s ease;
  }
  .toggle:hover:not(:disabled) .knob {
    background: radial-gradient(circle at 32% 28%, #f9fbfe, #dde4ed 55%, #b9c3d1);
  }
  .toggle.on .knob {
    transform: translateX(22px);
  }
  .toggle:active:not(:disabled) .knob {
    transform: translateY(0.5px) scale(0.95);
    box-shadow:
      inset 1px 1px 1px rgba(255, 255, 255, 0.4),
      inset -2px -3px 4px rgba(0, 0, 0, 0.18),
      1px 2px 3px rgba(0, 0, 0, 0.35);
  }
  .toggle.on:active:not(:disabled) .knob {
    transform: translateX(22px) translateY(0.5px) scale(0.95);
  }
</style>
