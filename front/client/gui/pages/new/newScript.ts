import api from '../../../api';

const empty = Object.create(null)

const elems = {
	get formElem() {
		const formElem = document.getElementById("new_message")
		if (!formElem) {
			console.error("[formElem]: no \"new_message\" form")
			return empty
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

	user = await api.getMe()
	if (user.error.error) {
		console.log("[onFormSubmit]: error loading me", user)
	}

	response = await api.uploadFile(elems.formElem.file.files[0])
	if (response.error.error) {
		console.log("[onFormSubmit]: error uploading file", response)
		return
	}

	response = await api.sendMessage(elems.formElem.text.value, response.id)
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