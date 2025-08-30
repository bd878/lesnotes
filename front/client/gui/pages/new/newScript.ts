import api from '../../../api';

const empty = Object.create(null)

const elems = {
	get formElem() {
		const formElem = document.getElementById("new_message")
		if (!formElem) {
			console.error("[formEleme]: no \"new_message\" form")
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

	console.log("submitted")
	await api.sendMessage("test form")
}