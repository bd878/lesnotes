import api from '../../../api';

const elems = {
	get formElem(): HTMLFormElement {
		const formElem = document.getElementById("new_message")
		if (!formElem) {
			console.error("[formElem]: no \"new_message\" form")
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	}
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

	let response, user;

	console.log("[onFormSubmit]: submitting...")

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

	if (user.error.error) {
		setTimeout(() => { location.href = "/m/" + response.message.ID }, 0)
	} else {
		setTimeout(() => { location.href = "/m/" + user.ID + "/" + response.message.ID }, 0)
	}
}