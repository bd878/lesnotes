import * as is from '../third_party/is'
import api from '../api'

async function publishMessage(ctx, next) {
	console.log("--> publishMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const redirectUrl = form.redirectUrl

	const response = await api.publishMessageJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		if (is.notEmpty(redirectUrl)) {
			ctx.redirect(redirectUrl)
		} else {
			ctx.redirect(ctx.router.url('message', {idOrName: form.id}, {query: ctx.query}))
		}
	}

	console.log("<-- publishMessage")
}

export default publishMessage;
