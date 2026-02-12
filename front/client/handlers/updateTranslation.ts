import updateTranslationJson from '../api/updateTranslationJson'
import * as is from '../third_party/is'

async function updateTranslation(ctx) {
	console.log("--> updateTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0

	const response = await updateTranslationJson(ctx.state.token, messageID, form.lang, form.title, form.text)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url("translation", {id: messageID, lang: form.lang}, {query: ctx.query}))
	}

	console.log("<-- updateTranslation")
}

export default updateTranslation
