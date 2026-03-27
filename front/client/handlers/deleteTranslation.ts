import deleteTranslationJson from '../api/deleteTranslationJson'
import * as is from '../third_party/is'

async function deleteTranslation(ctx) {
	console.log("--> deleteTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0;
	const redirectUrl = form.redirectUrl

	const response = await deleteTranslationJson(ctx.state.token, messageID, form.lang)

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

	console.log("<-- deleteTranslation")
}

export default deleteTranslation;
