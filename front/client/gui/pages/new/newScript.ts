import api from '../../../api';

const emptyElem = document.createElement("div")

const elems = {
	get formElem() {
		const formElem = document.getElementById("new_message")
		if (!formElem) {
			console.error("[formElem]: no \"new_message\" form")
			return emptyElem
		}

		return formElem
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
	if (user.error.error) {
		console.log("[onFormSubmit]: error loading me", user)
	}

	response = await api.uploadFile(elems.formElem.file.files[0])
	console.log("[onFormSubmit]: file:", response)
	if (response.error.error) {
		console.log("[onFormSubmit]: error uploading file", response)
		return
	}

	response = await api.sendMessage(elems.formElem.text.value, response.id)
	console.log("[onFormSubmit]: message:", response)
	if (response.error.error) {
		console.log("[onFormSubmit]: error saving message", response)
		return
	}

	if (user.error.error) {
		setTimeout(() => { location.href = "/m/" + response.message.ID }, 0)
	} else {
		setTimeout(() => { location.href = "/m/" + user.ID + "/" + response.message.ID }, 0)
	}
}