<script lang="ts">
  type props = {
    isLong?: boolean;
    disabled: boolean;
    push: (isLong: boolean) => Promise<void>;
  };

  const { isLong, disabled, push }: props = $props();
</script>

{#if isLong}
  <button
    type="button"
    class="long-push"
    {disabled}
    onclick={async () => {
      await push(true);
    }}>長押し</button
  >
{:else}
  <button
    type="button"
    class="push"
    {disabled}
    onclick={async () => {
      await push(false);
    }}>押す</button
  >
{/if}

<style>
  button {
    background: var(--background-gray);
    border: 2px solid var(--font-gray);
    padding: 0.25em 1em;
    border-radius: 0.375em;
    transition:
      border 250ms ease-out,
      background 250ms ease;
  }

  button:hover {
    border: 2px solid var(--primary-color);
    transition: border 100ms ease-in;
  }

  button:disabled {
    background: #ddd;
  }

  .push {
    font-size: 2rem;
  }
  .long-push {
    font-size: 1rem;
  }
</style>
