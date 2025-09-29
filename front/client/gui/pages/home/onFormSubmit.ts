import api from '../../../api';
import * as is from '../../../third_party/is';

async function onFormSubmit(elems, e) {
	e.preventDefault()

	if (either(elems.messageFormElem.messageText, elems.filesInputElem.files.length > 0)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}
	const user = await api.getMe()

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(params.get("cwd")) || 0

	const fileIDs = []

	if (elems.filesInputElem.files && is.notUndef(elems.filesInputElem.files[0])) {
		for (const file of elems.filesInputElem.files) {
			const response = await api.uploadFile(file)
			if (response.error.error) {
				console.error("[onFormSubmit]: cannot upload file:", response)
				return
			}

			fileIDs.push(response.ID)
		}
	}

	if (elems.messageFormElem.messageText) {
		const response = await api.sendMessage(elems.messageFormElem.messageText.value, elems.messageFormElem.messageTitle.value, fileIDs, threadID)
		if (response.error.error) {
			console.log("[onFormSubmit]: cannod send message:", response)
			return
		}
	}

	elems.messageFormElem.reset()

	location.href = params.toString() ? ("/home?" + params.toString()) : "/home"
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}

export default onFormSubmit
