import createTgAuth from '../../scripts/createTgAuth';
import onFormSubmit from './onFormSubmit';
import {getByID} from '../../../utils'

const elems = {
	form:   document.createElement("form"),
	div:    document.createElement("div"),

	get formElem():          HTMLFormElement    { return getByID("login-form",            this.form) as HTMLFormElement },
	get widgetElem():        HTMLDivElement     { return getByID("telegram-login-widget", this.div) as HTMLDivElement },
	get errorElem():         HTMLDivElement     { return getByID("login-error",           this.div) as HTMLDivElement },
}

function init() {
	console.log("loaded")

	elems.widgetElem.appendChild(createTgAuth())

	elems.formElem.addEventListener("submit", e => onFormSubmit(elems, e))
}

window.addEventListener("load", init)
