import * as is from '../third_party/is'
import api from '../api'

async function publishMessage(ctx, next) {
	console.log("--> publishMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const redirectUrl = form.redirectUrl
	const id = parseInt(form.id) || 0

	let response = await api.publishMessageJson(ctx.state.token, id)
	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	response = await api.publishThreadJson(ctx.state.token, id)
	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	if (is.notEmpty(redirectUrl)) {
		ctx.redirect(redirectUrl)
	} else {
		ctx.redirect(ctx.router.url('message', {idOrName: form.id}, {query: ctx.query}))
	}

	console.log("<-- publishMessage")
}

export default publishMessage;
