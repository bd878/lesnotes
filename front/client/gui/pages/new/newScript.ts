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

	let response;

	response = await api.uploadFile(elems.formElem.file[0])
	if (response.error.error) {
		console.log("[onFormSubmit]: error uploading file", response)
		return
	}

	response = await api.sendMessage(elems.formElem.text.value, response.id)
	if (response.error.error) {
		console.log("[onFormSubmit]: error saving message", response)
		return
	}

	setTimeout(() => { location.href = "/m/" + response.message.id }, 0)
}