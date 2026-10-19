import * as is from '../../../third_party/is'
import uploadFile from '../../../api/uploadFile';

async function onEditMessageFormSubmit(elems, e) {
	e.preventDefault()

	let name = ""
	if (is.notUndef(elems.messageEditFormElem.name)) {
		name = elems.messageEditFormElem.name.value
	}

	const fileIDs = new Set()

	const savedFiles = new Set()
	const nodes = elems.filesListElem.children
	for (let i = 0; i < nodes.length; i++) {
		savedFiles.add(nodes[i].dataset.name)
		const id = parseInt(nodes[i].id)
		if (!isNaN(id)) {
			fileIDs.add(id)
		}
	}

	if (elems.filesInputElem.files && is.notUndef(elems.filesInputElem.files[0])) {
		for (const file of elems.filesInputElem.files) {
			if (savedFiles.has(file.name)) {
				const response = await uploadFile(file)
				if (response.error.error) {
					console.error("[onEditMessageFormSubmit]: cannot upload file:", response)
					return
				}

				fileIDs.add(response.ID)
			}
		}
	}

	elems.messageEditFormElem.file_ids.value = JSON.stringify(Array.from(fileIDs))
	elems.messageEditFormElem.submit()
}

export default onEditMessageFormSubmit
