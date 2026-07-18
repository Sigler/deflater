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
    background: transparent;
    position: relative;
    transition:
      background 0.15s ease,
      border-color 0.15s ease;
  }
  .toggle:hover:not(:disabled) {
    border-color: var(--coral);
  }
  .toggle.on {
    background: var(--coral);
    border-color: var(--coral);
  }
  .toggle.on:hover:not(:disabled) {
    background: var(--coral-bright);
    border-color: var(--coral-bright);
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
    background: var(--text-dim);
    transition:
      transform 0.15s ease,
      background 0.15s ease;
  }
  .toggle.on .knob {
    transform: translateX(22px);
    background: var(--bg-window);
  }
</style>
