function editMessage(messageID) {
	const params = new URLSearchParams(location.search)
	params.set("edit", "1")
	params.set("id", messageID)

	location.href = "/home?" + params.toString()
}

export {
	editMessage,
}