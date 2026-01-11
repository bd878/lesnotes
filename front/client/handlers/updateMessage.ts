import updateMessageJson from '../api/updateMessageJson'
import * as is from '../third_party/is'

async function updateMessage(ctx) {
	console.log("--> updateMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const fileIDs = (form.file_ids || []).map(id => parseInt(id) || 0).filter(is.notEmpty)

	const response = await updateMessageJson(ctx.state.token, parseInt(form.id) || 0, form.text, form.title, form.name, fileIDs)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url('message', {id: form.id}, {query: ctx.query}))
	}

	console.log("<-- updateMessage")
}

export default updateMessage;
