import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function privateThread(ctx) {
	console.log("--> privateThread")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.privateThreadJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		ctx.redirect(ctx.router.url('thread', {id: form.id}, {query: ctx.query}))
	}

	console.log("<-- privateThread")
}

export default privateThread;
