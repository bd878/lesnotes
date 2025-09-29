import api from '../../../api';

async function onMessagePrivateClick(elems, e) {
	e.stopPropagation()
	const messageID = parseInt(elems.messagePrivateElem.dataset.messageId) || 0

	const response = await api.privateMessages([messageID])
	if (response.error.error) {
		console.error("[onMessagePrivateClick]: cannot private message:", response)
		return
	}

	location.reload()
}

export default onMessagePrivateClick
