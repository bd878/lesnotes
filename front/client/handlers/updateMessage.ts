import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function updateMessage(ctx) {
	console.log("--> update message")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.updateMessageJson(ctx.state.token, parseInt(form.id) || 0, form.text, form.title, form.name, [])

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		ctx.redirect(ctx.router.url('message', {id: form.id}, {query: ctx.query}))
	}

	console.log("<-- update message")
}

export default updateMessage;
