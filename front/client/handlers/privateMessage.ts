import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function privateMessage(ctx) {
	console.log("--> privateMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.privateMessageJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		ctx.redirect(ctx.router.url('message', {id: form.id}, {query: ctx.query}))
	}

	console.log("<-- privateMessage")
}

export default privateMessage;
