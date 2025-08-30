function init() {
	const scriptEl = document.createElement("script");
	scriptEl.setAttribute("async", true)
	scriptEl.setAttribute("src", "https://telegram.org/js/telegram-widget.js?22")
	scriptEl.setAttribute("data-telegram-login", `${BOT_USERNAME}`)
	scriptEl.setAttribute("data-size", "small")
	scriptEl.setAttribute("data-auth-url", `https://${BACKEND_URL}/tg_auth`)
	scriptEl.setAttribute("data-request-access", "write")
	const widgetElem = document.getElementById("telegram-login-widget")
	if (!widgetElem) {
		console.error("[mainScript]: no widget element")
		return
	}
	widgetElem.appendChild(scriptEl)
}

window.addEventListener("load", () => {
	console.log("loaded")
	init()
})
