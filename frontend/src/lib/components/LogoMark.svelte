<script lang="ts">
  // The Deflater mark: the four-pane window with one pane deflated,
  // slumping out of the grid with a little air escaping.
  let { size = 40 }: { size?: number } = $props();
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 96 96"
  fill="none"
  xmlns="http://www.w3.org/2000/svg"
  aria-hidden="true"
>
  <!-- three healthy panes -->
  <rect x="10" y="10" width="34" height="34" rx="7" fill="var(--coral)" />
  <rect x="52" y="10" width="34" height="34" rx="7" fill="var(--coral)" />
  <rect x="10" y="52" width="34" height="34" rx="7" fill="var(--coral)" />
  <!-- the fourth pane, out of air -->
  <path
    d="M54 55
       C54 53.3 55.3 52 57 52 H81
       C84 52 86.6 54.6 85.8 57.5
       C84.2 63.5 82.2 70.5 78.8 76.8
       C75.2 83.5 69.5 89.5 63.5 90.6
       C59.2 91.4 56.2 88.2 55.2 82.5
       C54.4 77.6 54 68 54 55 Z"
    fill="var(--coral)"
    opacity="0.62"
  />
  <!-- wrinkle -->
  <path
    d="M60 60 C61 70 62 78 65 84"
    stroke="var(--bg-window)"
    stroke-opacity="0.35"
    stroke-width="3"
    stroke-linecap="round"
  />
  <!-- escaping air -->
  <path
    d="M87 68 q5 -2.5 6.5 -7.5"
    stroke="var(--gold)"
    stroke-width="3.5"
    stroke-linecap="round"
  />
  <path
    d="M90 77 q6 -3 8 -9"
    stroke="var(--gold)"
    stroke-width="3.5"
    stroke-linecap="round"
    opacity="0.55"
  />
</svg>
