import type { Error } from '../../../api/models'
import api from '../../../api';

async function onFormSubmit(elems, e) {
	e.preventDefault()

	if (!elems.formElem.login) {
		console.error("[onFormSubmit]: form \"login-form\" has no field \"login\"")
		return
	}

	if (!elems.formElem.password) {
		console.error("[onFormSubmit]: form \"login-form\" has no field \"password\"")
		return
	}

	let login = elems.formElem.login.value
	let password = elems.formElem.password.value

	console.log("[onFormSubmit]: submitting", "login:", login, "password:", password)

	let response = await api.login(login, password)
	console.log("[onFormSubmit]: login:", response)
	if (response.error.error) {
		showError(elems, response.error)
		return
	}

	setTimeout(() => { location.href = "/home" }, 0)
}

function showError(elems: any, error: Error) {
	elems.errorElem.textContent = error.human
}

export default onFormSubmit
