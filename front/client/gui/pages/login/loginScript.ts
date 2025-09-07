import createTgAuth from '../../scripts/createTgAuth';
import api from '../../../api';

const elems = {
	get formElem(): HTMLFormElement {
		const formElem = document.getElementById("login-form")
		if (!formElem) {
			console.error("[loginScript]: no \"login-form\" form")
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	},

	get widgetElem() {
		const widgetElem = document.getElementById("telegram-login-widget")
		if (!widgetElem) {
			console.error("[loginScript]: no widget element")
			return document.createElement("div")
		}

		return widgetElem
	}
}

function init() {
	elems.widgetElem.appendChild(createTgAuth())
	elems.formElem.addEventListener("submit", onFormSubmit)
}

window.addEventListener("load", () => {
	console.log("loaded")
	init()
})

async function onFormSubmit(e) {
	e.preventDefault()

	if (!elems.formElem.login) {
		console.error("[onFormSubmit]: form \"login-form\" has no field \"login\"")
		return
	}

	if (!elems.formElem.password) {
		console.error("[onFormSubmit]: form \"login-form\" has no field \"password\"")
		return
	}

	// TODO: validate form
	// - check for input errors
	// - show error under form field
	// - set loading state

	let login = elems.formElem.login.value
	let password = elems.formElem.password.value

	console.log("[onFormSubmit]: submitting", "login:", login, "password:", password)

	let response = await api.login(login, password)
	console.log("[onFormSubmit]: login:", response)
	if (response.error.error) {
		// TODO: show form error
		return
	}

	setTimeout(() => { location.href = "/home" }, 0)
}