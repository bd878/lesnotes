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

	get messageEditFormElem(): HTMLFormElement {
		const formElem = document.getElementById("message-edit-form")
		if (!formElem) {
			console.error("[homeScript]: no \"message-edit-form\" form")
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
	},

	get threadsListElem(): HTMLDivElement {
		const divElem = document.getElementById("threads-list")
		if (!divElem) {
			console.error("[homeScript]: no \"threads-list\" elem")
			return document.createElement("div")
		}

		return divElem as HTMLDivElement
	},

	get messageDeleteElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-delete")
		if (!buttonElem) {
			console.error("[homeScript]: no \"message-delete\" elem")
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messageEditElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-edit")
		if (!buttonElem) {
			console.error("[homeScript]: no \"message-edit\" elem")
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messagePublishElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-publish")
		if (!buttonElem) {
			console.error("[homeScript]: no \"message-publish\" elem")
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messagePrivateElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-private")
		if (!buttonElem) {
			console.error("[homeScript]: no \"message-private\" elem")
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},

	get messageCancelEditElem(): HTMLButtonElement {
		const buttonElem = document.getElementById("message-cancel-edit")
		if (!buttonElem) {
			console.error("[homeScript]: no \"message-cancel-edit\" elem")
			return document.createElement("button")
		}

		return buttonElem as HTMLButtonElement
	},
}

function init() {
	elems.formElem.addEventListener("submit", onFormSubmit)
	elems.messageEditFormElem.addEventListener("submit", onMessageUpdateFormSubmit)
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

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

async function onMessageUpdateFormSubmit(e) {
	e.preventDefault()
}

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

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}