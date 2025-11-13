<script lang="ts">
  import type { NetworkStatusResponse } from "./types";
  import { fetchAPI, startStream } from "./fetch";

  let pulse = $state<boolean>(false);
  let isActive = $state<boolean>(false);

  let networkStatus = $state<NetworkStatusResponse>({
    ifaces: [],
    updated_at: "",
    error: false
  });

  async function updateNetworkStatus() {
    const result = await fetchAPI("/network/status", {
      method: "GET"
    });
    if (result instanceof Error) {
      networkStatus = {
        ifaces: [],
        updated_at: "",
        error: true
      };
      return;
    }

    try {
      const payload = JSON.parse(result);
      networkStatus = payload as NetworkStatusResponse;
    } catch (e) {
      networkStatus = {
        ifaces: [],
        updated_at: "",
        error: true
      };
    }
    resetCounter();
  }

  function getOffsetString(dt1: Date, dt2: Date) {
    const offsetTs = dt1.getTime() - dt2.getTime();

    if (offsetTs < 1 * 1000) {
      return "たった今";
    } else if (offsetTs < 60 * 1000) {
      return `${Math.floor(offsetTs / 1000)}秒前`;
    } else if (offsetTs < 60 * 60 * 1000) {
      return `${Math.floor(offsetTs / (60 * 1000))}分前`;
    } else if (offsetTs < 24 * 60 * 60 * 1000) {
      return `${Math.floor(offsetTs / (60 * 60 * 1000))}時間前`;
    } else {
      return `${Math.floor(offsetTs / (24 * 60 * 60 * 1000))}日前`;
    }
  }

  let latestUpdate = $state("");

  let fetched_at = new Date();
  const intervals: number[] = [];
  function resetCounter() {
    intervals.forEach((id) => clearInterval(id));
    fetched_at = new Date();
    const interval = setInterval(() => {
      latestUpdate = getOffsetString(new Date(), fetched_at);
    }, 5000);
    intervals.push(interval);

    latestUpdate = getOffsetString(new Date(), fetched_at);
  }

  (async () => {
    startStream(
      "/network/watch/stream",
      1500,
      [2500, 3000, 5000, 10000],
      (path, retryTimes) => {
        console.log(
          `Starting stream from: ${path} (retrying ${retryTimes} times)`
        );
        pulse = false;
        requestAnimationFrame(() => {
          pulse = true;
        });
      },
      (event) => {
        isActive = true;
        console.log("Stream opened:", event);
      },
      (data) => {
        try {
          const payload = JSON.parse(data);
          networkStatus = payload as NetworkStatusResponse;
          resetCounter();
        } catch {
          // 握りつぶす
        }
      },
      () => {
        pulse = false;
        requestAnimationFrame(() => {
          pulse = true;
        });
        resetCounter();
      },
      (message) => {
        isActive = false;
        console.log("Timeout!:", message);
      }
    );

    await updateNetworkStatus();
  })();
</script>

<section>
  <h2>ネットワークの状態</h2>
  <div>
    <div class="status">
      <!-- Switch element itself to trigger animation correctly -->
      <!-- TODO: a11y -->
      {#if pulse}
        <div class={["pulser", "pulse-animation", !isActive && "dead"]}></div>
      {:else}
        <div class={["pulser", !isActive && "dead"]}></div>
      {/if}
      <p>最終更新: {latestUpdate}</p>
      <button onclick={updateNetworkStatus}>今すぐ更新</button>
    </div>
    {#if networkStatus.error}
      <p class="error">ネットワークの状態を取得できませんでした</p>
    {/if}
    {#if networkStatus.ifaces.length > 0}
      <ul>
        {#each networkStatus.ifaces as iface}
          <li>
            <span class="name">{iface.name}</span>
            <span class={["status-icon", `status-icon_${iface.status}`]}></span>
            <span class="status-text">
              {#if iface.status == 0}
                未接続
              {:else if iface.status == 1}
                接続中
              {:else if iface.status == -1}
                不明
              {/if}
            </span>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</section>

<style>
  .pulse-animation {
    animation: pulse 1s;
  }
  @keyframes pulse {
    0% {
      box-shadow: 0 0 0 0 rgba(0, 0, 0, 0.2);
    }
    100% {
      box-shadow: 0 0 0 0.75rem rgba(0, 0, 0, 0);
    }
  }

  .status {
    display: flex;
    align-items: center;

    button {
      font-size: 0.875rem;
      background: var(--background-gray);
      border: 2px solid var(--font-gray);
      padding: 0.25em 0.5em;
      border-radius: 0.375em;
      transition:
        border 250ms ease-out,
        background 250ms ease;
      margin: 0 0.5em;
    }

    button:hover {
      border: 2px solid var(--primary-color);
      transition: border 100ms ease-in;
    }

    .pulser {
      width: 1.5rem;
      height: 1.5rem;
      border-radius: 50%;
      box-shadow: 0px 0px 1px 1px #0000001a;
      background-color: orange;
      margin: 1rem;
      &.dead {
        background-color: var(--font-gray);
      }
    }
  }
  .error {
    color: red;
  }

  ul {
    display: flex;
    list-style: none;
    padding: 0;
    gap: 1rem;
  }

  li {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-block: 0.5rem;
    border: 1px solid var(--font-gray);
    padding: 0.5rem;

    span {
      display: block;
    }

    .name {
      font-size: 1.2rem;
    }

    .status-text {
      font-size: 0.8rem;
    }

    .status-icon {
      width: 1rem;
      height: 1rem;
      border-radius: 50%;
      margin-block: 0.2rem;
    }

    .status-icon_0 {
      background-color: #dd0000;
    }

    .status-icon_1 {
      background-color: #00dd00;
    }

    .status-icon_-1 {
      background-color: var(--font-gray);
    }
  }
</style>
