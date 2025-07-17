<script lang="ts">
  import type { SpinnerStatuses } from "./types";

  type props = {
    status: SpinnerStatuses;
  };

  const { status }: props = $props();
</script>

<svg
  class={[
    "spinner",
    status === "loading" && "animate",
    status === "done" && "check"
  ]}
  viewBox="0 0 40 40"
  height="40"
  width="40"
  fill="none"
>
  <polyline
    class="check"
    points="5,22.5 17.5,32.5 35,10"
    stroke-width="6px"
    stroke-linejoin="round"
    stroke-linecap="round"
    stroke="#10a64a"
  />
  <circle
    class="track"
    cx="20"
    cy="20"
    r="17.5"
    pathLength="100"
    stroke-width="5px"
  />
  <circle
    class="car"
    cx="20"
    cy="20"
    r="17.5"
    pathLength="100"
    stroke-width="5px"
  />
</svg>

<style>
  .spinner {
    --uib-size: 1.5rem;
    --uib-color: black;
    --uib-speed: 2s;
    --uib-bg-opacity: 0;
    height: var(--uib-size);
    width: var(--uib-size);
    transform-origin: center;
    will-change: transform;
    overflow: visible;
    margin-inline: 0.5rem;

    .check {
      display: none;
    }
    .car,
    .track {
      display: none;
      fill: none;
      stroke: var(--uib-color);
      transition: stroke 0.5s ease;
    }

    .car {
      stroke-dasharray: 1, 200;
      stroke-dashoffset: 0;
      stroke-linecap: round;
      will-change: stroke-dasharray, stroke-dashoffset;
    }

    .track {
      opacity: var(--uib-bg-opacity);
    }
  }

  .spinner.animate {
    animation: rotate var(--uib-speed) linear infinite;

    .car,
    .track {
      display: unset;
    }
    .car {
      animation: stretch calc(var(--uib-speed) * 0.75) ease-in-out infinite;
    }
  }

  .spinner.check .check {
    display: unset;
  }

  @keyframes rotate {
    100% {
      transform: rotate(360deg);
    }
  }

  @keyframes stretch {
    0% {
      stroke-dasharray: 0, 150;
      stroke-dashoffset: 0;
    }
    50% {
      stroke-dasharray: 75, 150;
      stroke-dashoffset: -25;
    }
    100% {
      stroke-dashoffset: -100;
    }
  }
</style>
