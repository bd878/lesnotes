import * as is from '../third_party/is'
import api from '../api'

async function sendTranslation(ctx) {
	console.log("--> sendTranslation")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0

	const response = await api.sendTranslationJson(ctx.state.token, messageID, form.lang, form.text, form.title)
	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
	} else {
		ctx.redirect(ctx.router.url("translation", {id: messageID, lang: form.lang}, {query: ctx.query}))
	}

	console.log("<-- sendTranslation")
}

export default sendTranslation;
