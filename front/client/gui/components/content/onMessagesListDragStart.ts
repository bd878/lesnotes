
function onMessagesListDragStart(elems, e) {
	let id = e.target.dataset.messageId

	let rect = e.target.getBoundingClientRect()
	let shiftX = e.clientX - rect.left;
	let shiftY = e.clientY - rect.top;

	e.dataTransfer.effectAllowed = "move"
	e.dataTransfer.setData("text/plain", id)
	e.dataTransfer.setDragImage(e.target, shiftX, shiftY)
}

export default onMessagesListDragStart
