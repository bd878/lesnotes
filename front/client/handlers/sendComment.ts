import * as is from '../third_party/is'
import api from '../api'

async function sendComment(ctx) {
	console.log("--> sendComment")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0
	const redirectUrl = form.redirectUrl

	const response = await api.sendCommentJson(ctx.state.token, messageID, form.text)

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	} else {
		if (is.notEmpty(redirectUrl)) {
			ctx.redirect(redirectUrl)
		} else {
			ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
		}
	}

	console.log("<-- sendComment")
}

export default sendComment
