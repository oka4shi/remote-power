/** @type {HTMLButtonElement | null}*/
const pushButton = document.getElementById("push")
/** @type {HTMLButtonElement | null}*/
const longPushButton = document.getElementById("long-push")
const message = document.getElementById("message")
const spinner = document.getElementById("spinner")

/** @param {string} path
*   @param {RequestInit} init 
*/
async function fetchAPI(path, init) {
    let error = ""

    const res = await fetch(path, init).then((resp) => {
        if (!resp.ok) {
            if (resp.status == 409) {
                error = "現在実行中です！"
                return ""
            }
            error = "エラーが発生しました！"
            return ""
        }
        return resp.text()
    }).catch((err) => {
        error = `通信中にエラーが発生しました: ${err}`
        return ""
    })

    if (error) {
        return new Error(error)
    }

    return res
}

/** @param {boolean} isLong */
async function buttonClickHandler(isLong) {
    spinner.classList.remove("check")
    spinner.classList.add("animate")
    message.innerText = ``
    pushButton.disabled = true
    longPushButton.disabled = true

    const token = await fetchAPI(`/push${isLong ? "?long":""}`, {
        method: "POST"
    })
    if (token instanceof Error) {
        message.innerText = token.message
        message.classList.add("error")
        spinner.classList.remove("animate")
        pushButton.disabled = false
        longPushButton.disabled = false
        return
    }

    const result = await fetchAPI(`/push/status`, {
        method: "GET",
        headers: {
            "Push-Token": token
        }
    })
    if (result instanceof Error) {
        message.innerText = result.message
        message.classList.add("error")
        spinner.classList.remove("animate")
        pushButton.disabled = false
        longPushButton.disabled = false
        return
    }
    
    message.innerText = `操作が完了しました`
    message.classList.remove("error")
    spinner.classList.remove("animate")
    spinner.classList.add("check")
    pushButton.disabled = false
    longPushButton.disabled = false

}


pushButton?.addEventListener("click", async () => {await buttonClickHandler(false)})
longPushButton?.addEventListener("click", async () => {await buttonClickHandler(true)})
