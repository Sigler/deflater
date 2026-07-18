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
  .toggle {
    flex: none;
    width: 44px;
    height: 22px;
    border-radius: 999px;
    border: 1px solid var(--stroke-strong);
    background: rgba(0, 0, 0, 0.18);
    box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.3);
    position: relative;
    transition:
      background 0.15s ease,
      border-color 0.15s ease,
      box-shadow 0.15s ease;
  }
  .toggle:hover:not(:disabled) {
    border-color: var(--coral);
  }
  .toggle.on {
    background: linear-gradient(180deg, var(--coral-bright), var(--coral));
    border-color: var(--coral);
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.22),
      0 1px 2px rgba(0, 0, 0, 0.25);
  }
  .toggle.on:hover:not(:disabled) {
    background: linear-gradient(180deg, #ff8f70, var(--coral-bright));
  }
  .toggle:active:not(:disabled) {
    transform: translateY(0.5px);
  }
  .toggle:disabled {
    opacity: 0.45;
    cursor: default;
  }
  .knob {
    position: absolute;
    top: 3px;
    left: 4px;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: linear-gradient(180deg, #c7bcb1, var(--text-dim));
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
    transition:
      transform 0.15s ease,
      background 0.15s ease;
  }
  .toggle.on .knob {
    transform: translateX(22px);
    background: linear-gradient(180deg, #2e2723, var(--bg-window));
  }
</style>
