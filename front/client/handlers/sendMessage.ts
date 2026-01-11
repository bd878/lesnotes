import * as is from '../third_party/is'
import api from '../api'

async function sendMessage(ctx) {
	// TODO: proxy send message to messages service, /send
	console.log("--> sendMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	let fileIDs = []
	if (is.notEmpty(form.file_ids)) {
		if (is.array(form.file_ids)) {
			fileIDs = form.file_ids
		} else {
			fileIDs = [form.file_ids]
		}
	}

	fileIDs = fileIDs.map(id => parseInt(id) || 0).filter(is.notEmpty)

	const response = await api.sendMessageJson(ctx.state.token, form.text, form.title, fileIDs, parseInt(form.thread) || 0, true)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
	}

	console.log("<-- sendMessage")
}

export default sendMessage;
