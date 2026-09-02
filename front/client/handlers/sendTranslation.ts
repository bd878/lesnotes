import * as is from '../third_party/is'
import api from '../api'

async function sendTranslation(ctx) {
	ctx.log.info("--> sendTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0
	const redirectUrl = form.redirectUrl

	const response = await api.sendTranslationJson(ctx.state.token, messageID, form.lang, form.text, form.title)
	if (response.error.error) {
		ctx.log.error(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		if (is.notEmpty(redirectUrl)) {
			ctx.redirect(redirectUrl)
		} else {
			ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
		}
	}

	ctx.log.info("<-- sendTranslation")
}

export default sendTranslation;
