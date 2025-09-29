import * as is from '../../../third_party/is';
import api from '../../../api';

async function onMessageUpdateFormSubmit(elems, e) {
	e.preventDefault()

	if (is.notEmpty(e.target.dataset.messageId)) {
		const messageID = e.target.dataset.messageId

		const text = elems.editFormElem.messageText.value
		const title = elems.editFormElem.messageTitle.value

		let name = ""
		if (is.notUndef(elems.editFormElem.messageName)) {
			name = elems.editFormElem.messageName.value
		}

		const response = await api.updateMessage(messageID, text, title, name)
		if (response.error.error) {
			console.error("[onMessageUpdateFormSubmit]: cannot update message:", response)
			return
		}

		elems.editFormElem.reset()

		const params = new URL(location.toString()).searchParams
		params.delete("edit")

		location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
	} else {
		console.error("[onMessageUpdateFormSubmit]: no data-message-id attribute on target")
	}
}

export default onMessageUpdateFormSubmit
