import home from '../routes/home/home'
import * as is from '../third_party/is'
import api from '../api'

async function deleteMessage(ctx, next) {
	console.log("--> deleteMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const response = await api.deleteMessageJson(ctx.state.token, form.id)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		await home(ctx)
	} else {
		await next()
	}

	console.log("<-- deleteMessage")
}

export default deleteMessage;
