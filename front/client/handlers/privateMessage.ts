import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function privateMessage(ctx, next) {
	console.log("--> private message")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.privateMessageJson(ctx.state.token, form.id)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		await next()
	}

	console.log("<-- private message")
}

export default privateMessage;
