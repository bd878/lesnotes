import updateTranslationJson from '../api/updateTranslationJson'
import * as is from '../third_party/is'

async function updateTranslation(ctx) {
	console.log("--> updateTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0
	const redirectUrl = form.redirectUrl

	const response = await updateTranslationJson(ctx.state.token, messageID, form.lang, form.title, form.text)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		if (is.notEmpty(redirectUrl)) {
			ctx.redirect(redirectUrl)
		} else {
			ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
		}
	}

	console.log("<-- updateTranslation")
}

export default updateTranslation
