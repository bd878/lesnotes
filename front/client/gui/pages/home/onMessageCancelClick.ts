function onMessageCancelClick(elems, e) {
	e.stopPropagation()

	const params = new URLSearchParams(location.search)
	params.delete("id")

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

export default onMessageCancelClick
