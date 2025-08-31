import createTgAuth from '../../scripts/createTgAuth';

function init() {
	const widgetElem = document.getElementById("telegram-login-widget")
	if (!widgetElem) {
		console.error("[loginScript]: no widget element")
		return
	}
	widgetElem.appendChild(createTgAuth())
}

window.addEventListener("load", () => {
	console.log("loaded")
	init()
})
