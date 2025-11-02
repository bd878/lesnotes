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
		params.delete("cwd")
	} else {
		params.set("cwd", threadID)
	}

	params.delete("id")
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function paginateMessages(threadID, direction, offsetStr, limitStr) {
	const params = new URLSearchParams(location.search)

	const offset = parseInt(offsetStr)
	const limit = parseInt(limitStr)

	if (isNaN(offset) || isNaN(limit)) {
		console.error("[paginateMessages]: offset or limit are nan")
		return
	}

	if (direction == "prev") {
		params.set(threadID, `${offset + limit}`)
	} else if (direction == "next") {
		params.set(threadID, `${Math.max(0, offset - limit)}`)
	} else {
		console.error("[paginateMessages]: unknown direction:", direction)
	}

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

export {
	editMessage,
	showMessage,
	openThread,
	paginateMessages,
}