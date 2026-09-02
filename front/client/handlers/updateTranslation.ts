import updateTranslationJson from '../api/updateTranslationJson'
import * as is from '../third_party/is'

async function updateTranslation(ctx) {
	ctx.log.info("--> updateTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0
	const redirectUrl = form.redirectUrl

	const response = await updateTranslationJson(ctx.state.token, messageID, form.lang, form.title, form.text)

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

	ctx.log.info("<-- updateTranslation")
}

export default updateTranslation
