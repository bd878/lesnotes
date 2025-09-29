import api from '../../../api';

async function onMessageDeleteClick(elems, e) {
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

export default onMessageDeleteClick
