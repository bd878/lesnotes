import api from '../../../api';
import models from '../../../api/models';
import * as is from '../../../third_party/is';

const elems = {
	get formElem(): HTMLFormElement {
		const formElem = document.getElementById("message-form")
		if (!formElem) {
			console.error("[homeScript]: no \"message-form\" form")
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	},

	get messagesListElem(): HTMLDivElement {
		const divElem = document.getElementById("messages-list")
		if (!divElem) {
			console.error("[homeScript]: no \"messages-list\" elem")
			return document.createElement("div")
		}

		return divElem as HTMLDivElement
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

	if (either(elems.formElem.text, elems.formElem.file)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}
	const user = await api.getMe()

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(params.get("thread")) || 0

	if (elems.formElem.file && is.notUndef(elems.formElem.file.files[0])) {
		const file = await api.uploadFile(elems.formElem.file.files[0])
		if (file.error.error) {
			console.error("[onFormSubmit]: cannot upload file:", file)
			return
		}

		fileID = file.ID
	}

	if (elems.formElem.text) {
		const response = await api.sendMessage(elems.formElem.text.value, elems.formElem.messageTitle.value, fileID, threadID)
		if (response.error.error) {
			console.log("[onFormSubmit]: cannod send message:", response)
			return
		}
	}

	elems.formElem.reset()

	setTimeout(() => { location.href = "/home?" + params.toString() }, 0)
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}