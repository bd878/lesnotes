import * as is from '../third_party/is'
import api from '../api'

async function deleteMessage(ctx) {
	ctx.log.info("--> deleteMessage")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const redirectUrl = form.deleteRedirectUrl

	const response = await api.deleteMessageJson(ctx.state.token, parseInt(form.id) || 0)

	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		if (is.notEmpty(redirectUrl)) {
			ctx.redirect(redirectUrl)
		} else {
			ctx.redirect(ctx.router.url('home', {idOrName: form.id}, {query: ctx.query}))
		}
	}

	ctx.log.info("<-- deleteMessage")
}

export default deleteMessage;
