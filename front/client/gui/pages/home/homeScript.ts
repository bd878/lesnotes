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

	get fileFormElem(): HTMLFormElement {
		const formElem = document.getElementById("file-form")
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
	elems.formElem.addEventListener("submit", onFormSubmit)
	elems.filesButtonElem.addEventListener("click", onSelectFilesClick)
	elems.messageCancelElem.addEventListener("click", onMessageCancelClick)
	elems.editFormElem.addEventListener("submit", onMessageUpdateFormSubmit)
	elems.messagesListElem.addEventListener("click", onMessagesListClick)
	elems.threadsListElem.addEventListener("click", onThreadsListClick)
	elems.messageDeleteElem.addEventListener("click", onMessageDeleteClick)
	elems.messageEditElem.addEventListener("click", onMessageEditClick)
	elems.messagePublishElem.addEventListener("click", onMessagePublishClick)
	elems.messagePrivateElem.addEventListener("click", onMessagePrivateClick)
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

function onMessagesListClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		handleMessageClick(e.target.dataset.messageId)
	} else if (is.notEmpty(e.target.dataset.threadId)) {
		handleThreadClick(e.target.dataset.threadId)
	}
}

function onThreadsListClick(e) {
	if (is.notUndef(e.target.dataset.threadId)) {
		handleThreadClick(e.target.dataset.threadId)
	} else if (is.notEmpty(e.target.dataset.messageId)) {
		handleMessageClick(e.target.dataset.messageId)
	}
}

async function onMessagePublishClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = parseInt(e.target.dataset.messageId) || 0

		const response = await api.publishMessages([messageID])
		if (response.error.error) {
			console.error("[onMessagePublishClick]: cannot publish message:", response)
			return
		}

		location.reload()
	} else {
		console.error("[onMessagePublishClick]: no data-message-id attribute on target")
		return
	}
}

async function onMessagePrivateClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = parseInt(e.target.dataset.messageId) || 0

		const response = await api.privateMessages([messageID])
		if (response.error.error) {
			console.error("[onMessagePrivateClick]: cannot private message:", response)
			return
		}

		location.reload()
	} else {
		console.error("[onMessagePrivateClick]: no data-message-id attribute on target")
		return
	}
}

function onMessageCancelEditClick(e) {
	e.stopPropagation()

	const params = new URLSearchParams(location.search)
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

async function onMessageDeleteClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = parseInt(e.target.dataset.messageId) || 0

		const response = await api.deleteMessage(messageID)
		if (response.error.error) {
			console.error("[onMessageDeleteClick]: cannot delete message:", response)
			return
		}

		const params = new URLSearchParams(location.search)
		params.delete("id")

		location.href = params.toString() ? ("/home?" + params.toString()) : "/home" 
	} else {
		console.error("[onMessageDeleteClick]: no data-message-id attribute on target")
	}
}

function onMessageEditClick(e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = e.target.dataset.messageId

		const params = new URLSearchParams(location.search)
		params.set("edit", "1")

		location.href = "/home?" + params.toString()
	} else {
		console.error("[onMessageEditClick]: no data-message-id attribute on target")
	}
}

function handleMessageClick(messageID) {
	const params = new URLSearchParams(location.search)
	params.set("id", messageID)
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function handleThreadClick(threadID) {
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

async function onMessageUpdateFormSubmit(e) {
	e.preventDefault()

	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = e.target.dataset.messageId

		const text = elems.editFormElem.messageText.value
		const title = elems.editFormElem.messageTitle.value
		const name = elems.editFormElem.messageName.value

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

	if (either(elems.formElem.messageText, elems.fileFormElem.files.files.length > 0)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}
	const user = await api.getMe()

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(params.get("thread")) || 0

	const fileIDs = []

	if (elems.fileFormElem.files && is.notUndef(elems.fileFormElem.files.files[0])) {
		for (const file of elems.fileFormElem.files.files) {
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

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}