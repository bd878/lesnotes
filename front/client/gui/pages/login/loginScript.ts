import createTgAuth from '../../scripts/createTgAuth';
import getByID from '../../scripts/getByID'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),

	get widgetElem():            HTMLDivElement      { return getByID("telegram-login-widget", this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	elems.widgetElem.appendChild(createTgAuth())
}

window.addEventListener("load", init)
