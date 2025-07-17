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
