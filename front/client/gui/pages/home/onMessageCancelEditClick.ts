function onMessageCancelEditClick(elems, e) {
	e.stopPropagation()

	const params = new URLSearchParams(location.search)
	params.delete("edit")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

export default onMessageCancelEditClick