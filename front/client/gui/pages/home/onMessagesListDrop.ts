import * as is from '../../../third_party/is';
import api from '../../../api';

async function onMessagesListDrop(elems, e) {
	e.preventDefault()
	if (is.empty(e.target.dataset) || is.empty(e.target.dataset.messageId)) {
		console.log("[onMessagesListDrop]: no data-message-id")
		return
	}

	const targetId = e.target.dataset.messageId
	const targetElem = document.getElementById("list-" + targetId)

	if (is.undef(targetElem)) {
		console.log("[onMessagesListDrop]: cannot find target element by id", "list-" + targetId)
		return
	}

	if (is.undef(e.dataTransfer)) {
		console.log("[onMessagesListDrop]: no dataTransfer")
		return
	}

	const sourceId = e.dataTransfer.getData("text")
	const sourceElem = document.getElementById("list-" + sourceId)
	if (is.undef(sourceElem)) {
		console.log("[onMessagesListDrop]: cannot find source element by id", "list-" + sourceId)
		return
	}

	const response = await api.reorderThread(sourceId, -1 /* no parent */, targetId, 0 /* no prev */)
	if (response.error.error) {
		console.error('[onMessagesListDrop]: cannot reorder thread:', response)
		return
	}

	targetElem.insertAdjacentElement('beforebegin', sourceElem)
}

export default onMessagesListDrop;
