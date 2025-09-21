export async function fetchAPI(path: string, init: RequestInit) {
  let error = "";

  const res = await fetch(path, init)
    .then((resp) => {
      if (!resp.ok) {
        if (resp.status == 409) {
          error = "現在実行中です！";
          return "";
        }
        error = "エラーが発生しました！";
        return "";
      }
      return resp.text();
    })
    .catch((err) => {
      error = `通信中にエラーが発生しました: ${err}`;
      return "";
    });

  if (error) {
    return new Error(error);
  }

  return res;
}

export function startStream(
  path: string,
  timeout: number,
  retryIn: number,
  onStart: (path: string, retryTimes: number) => void,
  onOpen: (event: Event) => void,
  onMessage: (data: string) => void,
  onHeartbeat: () => void,
  onTimeout: (message: string) => void,
  retryTimes: number = 0
): EventSource | null {
  onStart(path, retryTimes);
  let errorTimeout: number | undefined = undefined;
  const eventSource = new EventSource(path);

  errorTimeout = setTimeout(
    () => {
      onTimeout(
        `Can't connect to stream in ${retryIn * (retryTimes * 0.5 + 1)} milliseconds`
      );
      stopStream(eventSource, errorTimeout);
      startStream(
        path,
        timeout,
        retryIn,
        onStart,
        onOpen,
        onMessage,
        onHeartbeat,
        onTimeout,
        retryTimes + 1
      );
    },
    retryIn * (retryTimes * 0.5 + 1)
  );

  eventSource.onopen = (event) => {
    clearTimeout(errorTimeout);
    retryTimes = 0;
    onOpen(event);
  };

  eventSource.onmessage = (event) => {
    clearTimeout(errorTimeout);
    onMessage(event.data);
  };

  eventSource.addEventListener("heartbeat", () => {
    clearTimeout(errorTimeout);
    onHeartbeat();

    errorTimeout = setTimeout(() => {
      onTimeout(`No heartbeat received in ${timeout} milliseconds`);
      stopStream(eventSource, errorTimeout);
      startStream(
        path,
        timeout,
        retryIn,
        onStart,
        onOpen,
        onMessage,
        onHeartbeat,
        onTimeout
      );
    }, timeout);
  });

  eventSource.onerror = () => {
    onTimeout("Stream error occurred");
    clearTimeout(errorTimeout);

    // Retry connection instantly
    stopStream(eventSource, errorTimeout);
    startStream(
      path,
      timeout,
      retryIn,
      onStart,
      onOpen,
      onMessage,
      onHeartbeat,
      onTimeout
    );
  };

  return eventSource;
}

function stopStream(
  eventSource: EventSource,
  errorTimeout: number | undefined
) {
  if (eventSource) {
    eventSource.close();
  }
  clearTimeout(errorTimeout);
}
