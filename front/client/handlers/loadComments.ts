import listCommentsJson, { EmptyListComments } from '../api/listCommentsJson'
import * as is from '../third_party/is';

const commentsLimit = 1_000

async function loadComments(ctx, next) {
	const id = ctx.state.messageID || 0
	const name = ctx.state.messageName || ""
	const token = ctx.state.token

	ctx.log.info("--> loadComments", "id", id, "name", name, "token", token)

	ctx.state.comments = await listCommentsJson(token, id, name, commentsLimit, 0)
	if (is.notEmpty(ctx.state.comments)) {
		if (ctx.state.comments.error.error) {
			ctx.log.error(ctx.state.comments.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.comments = ctx.state.comments.comments
	} else {
		ctx.state.comments = EmptyListComments.comments
	}

	ctx.state.comments.reverse()

	await next()

	ctx.log.info("<-- loadComments")
}

export default loadComments
