import updateMessage from '../../../api/updateMessage'
import * as is from '../../../third_party/is'

async function autosaveMessage(elems, _e) {
	const id = (elems.messageEditForm.elements.namedItem("id") || elems.input).value
	const name = (elems.messageEditForm.elements.namedItem("name") || elems.input).value
	const title = (elems.messageEditForm.elements.namedItem("title") || elems.input).value
	const text = (elems.messageEditForm.elements.namedItem("text") || elems.textarea).value
	const fileIDs = getSelectedValues((elems.messageEditForm.elements.namedItem("file_ids") || elems.select).selectedOptions)
		.map(id => parseInt(id) || 0).filter(is.notEmpty)

	console.log("autosaveMessage", "id", id)

	try {
		await updateMessage(id, text, title, name, fileIDs)
	} catch (e) {
		console.error(e)
	}
}

function getSelectedValues(options: HTMLOptionElement[]): number[] {
	let result = []
	for (let i = 0; i < options.length; result.push(options[i].value), i++) {}
	return result
}

export default autosaveMessage
