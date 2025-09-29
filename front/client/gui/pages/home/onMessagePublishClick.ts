import api from '../../../api';

async function onMessagePublishClick(elems, e) {
	e.stopPropagation()
	const messageID = parseInt(elems.messagePublishElem.dataset.messageId) || 0

	const response = await api.publishMessages([messageID])
	if (response.error.error) {
		console.error("[onMessagePublishClick]: cannot publish message:", response)
		return
	}

	location.reload()
}

export default onMessagePublishClick
