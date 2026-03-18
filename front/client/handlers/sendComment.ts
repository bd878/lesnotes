import * as is from '../third_party/is'
import api from '../api'

async function sendComment(ctx) {
	console.log("--> sendComment")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0

	let response;
	if (is.notEmpty(ctx.state.token)) {
		response = await api.sendCommentJson(ctx.state.token, messageID, form.text)
	} else {
		ctx.body = "unimplemented"
		ctx.status = 501
		return
	}

	if (response.error.error) {
		console.log(response.error)
		ctx.state.error = response.error.human
		ctx.body = "error"
		return
	} else {
		ctx.redirect(ctx.router.url('home', {}, {query: ctx.query}))
	}

	console.log("<-- sendComment")
}

export default sendComment
