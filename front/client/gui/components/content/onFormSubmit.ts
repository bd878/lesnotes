import uploadFile from '../../../api/uploadFile';
import * as is from '../../../third_party/is';

const limit = parseInt(LIMIT)

async function onNewMessageFormSubmit(elems, e) {
	e.preventDefault()

	if (either(elems.newMessageFormElem.text, elems.filesInputElem.files.length > 0)) {
		console.error("[onFormSubmit]: either text of file must be present")
		return
	}

	let fileID = 0;

	const params = new URL(location.toString()).searchParams
	const threadID = parseInt(elems.newMessageFormElem.thread.value) || 0

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

	elems.newMessageFormElem.file_ids.value = JSON.stringify(fileIDs)
	elems.newMessageFormElem.submit()
}

function either(st1: boolean, st2: boolean): boolean {
	return (!st1 && !st2)
}

export default onNewMessageFormSubmit