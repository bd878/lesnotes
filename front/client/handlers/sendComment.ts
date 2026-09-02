import * as is from '../third_party/is'
import api from '../api'

async function sendComment(ctx) {
	ctx.log.info("--> sendComment")

	let form = ctx.request.body

	if (is.empty(form)) {
		form = {}
	}

	const messageID = parseInt(form.message) || 0
	const redirectUrl = form.redirectUrl

	const response = await api.sendCommentJson(ctx.state.token, messageID, form.text)

	if (response.error.error) {
		ctx.log.error(response.error)
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

	ctx.log.info("<-- sendComment")
}

export default sendComment
