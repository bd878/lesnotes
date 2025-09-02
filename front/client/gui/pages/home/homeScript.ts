import api from '../../../api';

const elems = {
	get formElem(): HTMLFormElement {
		const formElem = document.getElementById("message-form")
		if (!formElem) {
			console.error("[homeScript]: no \"message-form\" form")
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	},
}

function init() {
	elems.formElem.addEventListener("submit", onFormSubmit)
}

window.addEventListener("load", () => {
	console.log("loaded")
	init()
})

async function onFormSubmit(e) {
	e.preventDefault()

	if (!elems.formElem.text) {
		console.error("[onFormSubmit]: form \"message-form\" has no field \"text\"")
		return
	}

	if (!elems.formElem.file) {
		console.error("[onFormSubmit]: form \"message-form\" has no field \"file\"")
		return
	}

	console.log("[onFormSubmit]: submitting...")

	let response, user;

	user = await api.getMe()
	console.log("[onFormSubmit]: user:", user)

	response = await api.uploadFile(elems.formElem.file.files[0])
	console.log("[onFormSubmit]: file:", response)
	if (response.error.error) {
		return
	}

	response = await api.sendMessage(elems.formElem.text.value, response.id)
	console.log("[onFormSubmit]: message:", response)
	if (response.error.error) {
		return
	}

	elems.formElem.reset()

	const params = new URL(location.toString()).searchParams
	setTimeout(() => { location.href = "/home" + params.toString() }, 0)
}