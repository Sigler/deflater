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
  /* The track is a recessed well, lit from above: the rim is drawn with
     an inset ring (no real border, which would skew the geometry and
     catch false light), the top inner lip is in shadow, and the only
     light is a faint reflection on the well floor and a hairline on the
     bottom outer lip. The state color sits down inside the well and
     stays matte; nothing on top of the shape shines. */
  .toggle {
    flex: none;
    width: 44px;
    height: 22px;
    border-radius: 999px;
    border: none;
    background: linear-gradient(180deg, rgba(0, 0, 0, 0.38), rgba(0, 0, 0, 0.18));
    box-shadow:
      inset 0 0 0 1px rgba(0, 0, 0, 0.4),
      inset 0 2px 3px rgba(0, 0, 0, 0.5),
      inset 0 -1px 1px rgba(255, 244, 230, 0.05),
      0 1px 0 rgba(255, 244, 230, 0.06);
    position: relative;
    transition: background 0.15s ease;
  }
  .toggle.on {
    background: linear-gradient(180deg, #8f3a2a, #cd5a3d 60%, #dd664a);
    box-shadow:
      inset 0 0 0 1px rgba(0, 0, 0, 0.42),
      inset 0 2px 3px rgba(0, 0, 0, 0.45),
      inset 0 -1px 1px rgba(255, 200, 170, 0.12),
      0 1px 0 rgba(255, 244, 230, 0.06);
  }
  .toggle.on:hover:not(:disabled) {
    background: linear-gradient(180deg, #9a4030, #d9603f 60%, #e77049);
  }
  .toggle:disabled {
    opacity: 0.45;
    cursor: default;
  }
  /* The knob is a raised button sitting proud of the well: lit from
     above, casting a soft shadow down into it. Hover brightens the
     button itself rather than ringing the well. */
  .knob {
    position: absolute;
    top: 3px;
    left: 3px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: linear-gradient(180deg, #f2eae0, #cfc3b6);
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.55),
      inset 0 -1px 1px rgba(0, 0, 0, 0.12),
      0 1px 1px rgba(0, 0, 0, 0.4),
      0 2px 3px rgba(0, 0, 0, 0.3);
    transition:
      transform 0.15s ease,
      background 0.15s ease;
  }
  .toggle:hover:not(:disabled) .knob {
    background: linear-gradient(180deg, #f9f3ea, #ddd1c4);
  }
  .toggle.on .knob {
    transform: translateX(22px);
  }
  .toggle:active:not(:disabled) .knob {
    transform: translateY(0.5px) scale(0.96);
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.4),
      0 1px 1px rgba(0, 0, 0, 0.35);
  }
  .toggle.on:active:not(:disabled) .knob {
    transform: translateX(22px) translateY(0.5px) scale(0.96);
  }
</style>
