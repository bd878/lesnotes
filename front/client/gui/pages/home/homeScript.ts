import api from '../../../api';
import models from '../../../api/models';
import * as is from '../../../third_party/is';

const elems = {
	get formElem(): HTMLFormElement {
		const formElem = document.getElementById("message-form")
		if (!formElem) {
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	},

	get filesButtonElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("select-files-button")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get filesListElem(): HTMLDivElement {
		const divElem = document.getElementById("files-list")
		if (!divElem) {
			return document.createElement("div")
		}

		return divElem as HTMLDivElement
	},

	get filesInputElem(): HTMLInputElement {
		const inputElem = document.getElementById("files-input")
		if (!inputElem) {
			return document.createElement("input")
		}

		return inputElem as HTMLInputElement
	},

	get editFormElem(): HTMLFormElement {
		const formElem = document.getElementById("message-edit-form")
		if (!formElem) {
			return document.createElement("form")
		}

		return formElem as HTMLFormElement
	},

	get messagesListElem(): HTMLDivElement {
		const divElem = document.getElementById("messages-list")
		if (!divElem) {
			return document.createElement("div")
		}

		return divElem as HTMLDivElement
	},

	get threadsListElem(): HTMLDivElement {
		const divElem = document.getElementById("threads-list")
		if (!divElem) {
			return document.createElement("div")
		}

		return divElem as HTMLDivElement
	},

	get messageDeleteElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-delete")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messageEditElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-edit")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messagePublishElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-publish")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messagePrivateElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-private")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messageCancelEditElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-cancel-edit")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messageCancelElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-cancel")
		if (!buttonElem) {
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	}
}

function init() {
	elems.formElem.addEventListener("submit",             onFormSubmit)
	elems.filesInputElem.addEventListener("change",       onFileInputChange)
	elems.filesButtonElem.addEventListener("click",       onSelectFilesClick)
	elems.messageCancelElem.addEventListener("click",     onMessageCancelClick)
	elems.editFormElem.addEventListener("submit",         onMessageUpdateFormSubmit)
	elems.messagesListElem.addEventListener("click",      onMessagesListClick)
	elems.threadsListElem.addEventListener("click",       onThreadsListClick)
	elems.messageDeleteElem.addEventListener("click",     onMessageDeleteClick)
	elems.messageEditElem.addEventListener("click",       onMessageEditClick)
	elems.messagePublishElem.addEventListener("click",    onMessagePublishClick)
	elems.messagePrivateElem.addEventListener("click",    onMessagePrivateClick)
	elems.messageCancelEditElem.addEventListener("click", onMessageCancelEditClick)
}

window.addEventListener("load", () => {
	console.log("loaded")
	init()
})

function onSelectFilesClick(e) {
	if (is.notEmpty(elems.filesInputElem.id)) {
		elems.filesInputElem.click()
	}
}

function onMessageCancelClick(e) {
	e.stopPropagation()

	const params = new URLSearchParams(location.search)
	params.delete("id")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function createFilesListElement(fileName: string): HTMLDivElement {
	const elem = document.createElement("div")

	const textElem = document.createElement("span")
	const removeButton = document.createElement("button")

	removeButton.textContent = "X"
	removeButton.classList.add(...("cursor-pointer underline hover:text-blue-600 mr-2").split(" "))

	removeButton.onclick = () => { elem.remove() }

	textElem.textContent = fileName
	textElem.classList.add(...("overflow-hidden text-ellipsis").split(" "))

	elem.classList.add(...("mb-2 overflow-hidden text-ellipsis".split(" ")))
	elem.appendChild(removeButton)
	elem.appendChild(textElem)

	return elem
}

function onFileInputChange(e) {
	for (const file of e.target.files) {
		elems.filesListElem.appendChild(createFilesListElement(file.name))
	}
}

function onMessagesListClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		showMessage(e.target.dataset.messageId)
	} else if (is.notEmpty(e.target.dataset.direction) && is.notEmpty(e.target.dataset.threadId)) {
		paginateMessages(e.target.dataset.threadId, e.target.dataset.direction)
	} else if (is.notEmpty(e.target.dataset.threadId)) {
		openThread(e.target.dataset.threadId)
	}
}

function onThreadsListClick(e) {
	if (is.notUndef(e.target.dataset.threadId)) {
		openThread(e.target.dataset.threadId)
	} else if (is.notEmpty(e.target.dataset.messageId)) {
		showMessage(e.target.dataset.messageId)
	}
}

function onMessageEditClick(e) {
	e.stopPropagation()
	editMessage(parseInt(elems.messageEditElem.dataset.messageId))
}

async function onMessagePublishClick(e) {
	e.stopPropagation()
	const messageID = parseInt(elems.messagePublishElem.dataset.messageId) || 0

	const response = await api.publishMessages([messageID])
	if (response.error.error) {
		console.error("[onMessagePublishClick]: cannot publish message:", response)
		return
	}

	location.reload()
}

async function onMessagePrivateClick(e) {
	e.stopPropagation()
	const messageID = parseInt(elems.messagePrivateElem.dataset.messageId) || 0

	const response = await api.privateMessages([messageID])
	if (response.error.error) {
		console.error("[onMessagePrivateClick]: cannot private message:", response)
		return
	}

	location.reload()
}

function onMessageCancelEditClick(e) {
	e.stopPropagation()

	const params = new URLSearchParams(location.search)
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

async function onMessageDeleteClick(e) {
	e.stopPropagation()
	const messageID = parseInt(elems.messageDeleteElem.dataset.messageId) || 0

	const response = await api.deleteMessage(messageID)
	if (response.error.error) {
		console.error("[onMessageDeleteClick]: cannot delete message:", response)
		return
	}

	const params = new URLSearchParams(location.search)
	params.delete("id")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home" 
}

async function onMessageUpdateFormSubmit(e) {
	e.preventDefault()

	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = e.target.dataset.messageId

		const text = elems.editFormElem.messageText.value
		const title = elems.editFormElem.messageTitle.value

		let name = ""
		if (is.notUndef(elems.editFormElem.messageName)) {
			name = elems.editFormElem.messageName.value
		}

		const response = await api.updateMessage(messageID, text, title, name)
		if (response.error.error) {
			console.error("[onMessageUpdateFormSubmit]: cannot update message:", response)
			return
		}

		elems.editFormElem.reset()

		const params = new URL(location.toString()).searchParams
		params.delete("edit")

		location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
	} else {
		console.error("[onMessageUpdateFormSubmit]: no data-message-id attribute on target")
	}
}

async function onFormSubmit(e) {
	e.preventDefault()

	if (either(elems.formElem.messageText, elems.filesInputElem.files.length > 0)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}
	const user = await api.getMe()

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(params.get("thread")) || 0

	const fileIDs = []

	if (elems.filesInputElem.files && is.notUndef(elems.filesInputElem.files[0])) {
		for (const file of elems.filesInputElem.files) {
			const response = await api.uploadFile(file)
			if (response.error.error) {
				console.error("[onFormSubmit]: cannot upload file:", response)
				return
			}

			fileIDs.push(response.ID)
		}
	}

	if (elems.formElem.messageText) {
		const response = await api.sendMessage(elems.formElem.messageText.value, elems.formElem.messageTitle.value, fileIDs, threadID)
		if (response.error.error) {
			console.log("[onFormSubmit]: cannod send message:", response)
			return
		}
	}

	elems.formElem.reset()

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function paginateMessages(threadID, direction) {/*TODO: implement*/}

function editMessage(messageID) {
	const params = new URLSearchParams(location.search)
	params.set("edit", "1")
	params.set("id", messageID)

	location.href = "/home?" + params.toString()
}

function showMessage(messageID) {
	const params = new URLSearchParams(location.search)
	params.set("id", messageID)
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function openThread(threadID) {
	const params = new URLSearchParams(location.search)
	if (threadID == 0 || threadID == "0" || threadID == "") {
		params.delete("thread")
	} else {
		params.set("thread", threadID)
	}

	params.delete("id")
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}