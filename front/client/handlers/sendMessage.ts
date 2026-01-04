import home from '../routes/home'
import * as is from '../third_party/is'
import api from '../api'

async function sendMessage(ctx) {
	// TODO: proxy send message to messages service, /send
	console.log("--> sendMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.sendMessageJson(ctx.state.token, form.text, form.title, form.file_ids, parseInt(form.thread) || 0, true)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
	}

	console.log("<-- sendMessage")
}

export default sendMessage;
