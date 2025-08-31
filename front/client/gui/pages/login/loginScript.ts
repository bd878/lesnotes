import createTgAuth from '../../scripts/createTgAuth';
import api from '../../../api';

const emptyElem = Document.createElement("div")

const elems = {
	get formElem() {
		const formElem = document.getElementById("login-form")
		if (!formElem) {
			console.error("[loginScript]: no \"login-form\" form")
			return emptyElem
		}

		return formElem
	}

	get widgetElem() {
		const widgetElem = document.getElementById("telegram-login-widget")
		if (!widgetElem) {
			console.error("[loginScript]: no widget element")
			return emptyElem
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

	if (!elems.formElem.name) {
		console.error("[onFormSubmit]: form \"login-form\" has no field \"name\"")
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

	let name = elems.formElem.name.value
	let password = elems.formElem.password.value

	console.log("[onFormSubmit]: submitting", "name:", name, "password:", password)

	let response = await api.login(name, password)
	console.log("[onFormSubmit]: login:", response)
	if (response.error.error) {
		console.error("[onFormSubmit]: error logging in", response)
		// TODO: show form error
		return
	}

	setTimeout(() => { location.href = "/home" }, 0)
}