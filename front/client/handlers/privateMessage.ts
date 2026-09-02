import * as is from '../third_party/is'
import api from '../api'

async function privateMessage(ctx) {
	ctx.log.info("--> privateMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const redirectUrl = form.redirectUrl
	const id = parseInt(form.id) || 0

	let response = await api.privateMessageJson(ctx.state.token, id)
	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	response = await api.privateThreadJson(ctx.state.token, id)
	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	}

	if (is.notEmpty(redirectUrl)) {
		ctx.redirect(redirectUrl)
	} else {
		ctx.redirect(ctx.router.url('message', {idOrName: form.id}, {query: ctx.query}))
	}

	ctx.log.info("<-- privateMessage")
}

export default privateMessage;
