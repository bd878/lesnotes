import getMe from '../../../api/getMe';
import uploadFile from '../../../api/uploadFile';
import sendMessage from '../../../api/sendMessage';
import * as is from '../../../third_party/is';

const limit = parseInt(LIMIT)

async function onFormSubmit(elems, e) {
	e.preventDefault()

	if (either(elems.newMessageFormElem.text, elems.filesInputElem.files.length > 0)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}
	const user = await getMe()

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(params.get("cwd")) || 0

	const fileIDs = []

	const savedFiles = new Set
	const nodes = elems.filesListElem.children
	for (let i = 0; i < nodes.length; i++) {
		savedFiles.add(nodes[i].dataset.name)
	}

	if (elems.filesInputElem.files && is.notUndef(elems.filesInputElem.files[0])) {
		for (const file of elems.filesInputElem.files) {
			if (savedFiles.has(file.name)) {
				const response = await uploadFile(file)
				if (response.error.error) {
					console.error("[onFormSubmit]: cannot upload file:", response)
					return
				}

				fileIDs.push(response.ID)
			}
		}
	}

	let response
	if (elems.newMessageFormElem.text) {
		response = await sendMessage(elems.newMessageFormElem.text.value, elems.newMessageFormElem.title.value, fileIDs, threadID)
		if (response.error.error) {
			console.log("[onFormSubmit]: cannot send message:", response)
			return
		}
	}

	elems.newMessageFormElem.reset()

	params.set(`${threadID}`, `${limit},0`)
	location.href = params.toString() ? ("/messages/" + response.message.ID + "?" + params.toString()) : "/home"
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}

export default onFormSubmit