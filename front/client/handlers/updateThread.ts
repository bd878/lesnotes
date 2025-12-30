import home from '../routes/home'
import * as is from '../third_party/is'
import api from '../api'

async function updateThread(ctx) {
	console.log("--> updateThread")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.updateThreadJson(ctx.state.token, parseInt(form.id) || 0, form.description, form.name)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		ctx.redirect(ctx.router.url('thread', {id: form.id}, {query: ctx.query}))
	}

	console.log("<-- updateThread")
}

export default updateThread;
