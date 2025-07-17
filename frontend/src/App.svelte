<script lang="ts">
  import PushButton from "./lib/PowerButton.svelte";
  import Spinner from "./lib/Spinner.svelte";
  import NetworkStatus from "./lib/NetworkStatus.svelte";

  import { fetchAPI } from "./lib/fetch";
  import type { SpinnerStatuses, NetworkStatusResponse } from "./lib/types";

  let status = $state<SpinnerStatuses>("none");
  let disabled = $state(false);
  let message = $state({
    message: "",
    isError: false
  });

  async function push(isLong: boolean) {
    status = "loading";
    message = {
      message: "",
      isError: false
    };
    disabled = true;

    const token = await fetchAPI(`/push${isLong ? "?long" : ""}`, {
      method: "POST"
    });
    if (token instanceof Error) {
      message = {
        message: token.message,
        isError: true
      };
      status = "none";
      disabled = false;
      return;
    }

    const result = await fetchAPI(`/push/status`, {
      method: "GET",
      headers: {
        "Push-Token": token
      }
    });
    if (result instanceof Error) {
      message = {
        message: result.message,
        isError: true
      };
      status = "none";
      disabled = false;
      return;
    }

    message = {
      message: "操作が完了しました",
      isError: false
    };
    status = "done";
    disabled = false;
  }
</script>

<main id="container">
  <h1>電源管理ツール</h1>
  <p>自宅のデスクトップPCの電源を管理します</p>

  <section>
    <h2>電源ボタンを押す</h2>
    <div id="control">
      <PushButton isLong={false} {push} {disabled} />
      <PushButton isLong={true} {push} {disabled} />
    </div>

    <div id="status">
      <Spinner {status} />
      <p class={[{ error: message.isError }]}>{message.message}</p>
    </div>
  </section>

  <NetworkStatus />
</main>

<style>
  #container {
    background-color: var(--background-white);
    border-radius: 0.5em;
    padding: 1rem;
    max-width: 1500px;
    min-height: calc(100dvh - 2rem);
    margin: 1rem auto;
  }
  #control {
    display: flex;
    align-items: center;
    column-gap: 1rem;
  }
  #status {
    display: flex;
    align-items: center;
    margin-block: 0.5rem;

    p {
      margin-block: 0;
    }

    .error {
      color: red;
    }
  }
</style>
