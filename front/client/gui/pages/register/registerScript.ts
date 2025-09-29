import createTgAuth from '../../scripts/createTgAuth';

function init() {
	console.log("loaded")

	const widgetElem = document.getElementById("telegram-login-widget")
	if (!widgetElem) {
		console.error("[registerScript]: no widget element")
		return
	}

	widgetElem.appendChild(createTgAuth())
}

window.addEventListener("load", init)
