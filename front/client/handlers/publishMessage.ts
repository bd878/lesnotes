import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function publishMessage(ctx, next) {
	console.log("--> publish message")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.publishMessageJson(ctx.state.token, form.id)

	if (response.error.error) {
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		await next()
	}

	console.log("<-- publish message")
}

export default publishMessage;
