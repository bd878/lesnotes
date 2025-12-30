import api from '../api';
import models from '../api/models';
import * as is from '../third_party/is';

const limit = parseInt(LIMIT)

async function loadThreadMessages(ctx, next) {
	const userID = is.notEmpty(ctx.state.thread) ? ctx.state.thread.userID : (parseInt(ctx.params.user) || 0)
	const id = is.notEmpty(ctx.state.thread) ? ctx.state.thread.ID : 0;
	const token = ctx.state.token
	const params = new URLSearchParams(ctx.search)
	const offset = parseInt(params.get("offset")) || 0

	console.log("--> loadThreadMessages")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.messages = await api.readMessagesJson(token, 0 /* me */, id, 1 /* order */, limit, offset)
		} else {
			// TODO: load private thread messages by name
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			// load public thread messages (userID is neccessary)
			ctx.state.messages = await api.readMessagesJson("", userID, id, 1 /* order */, limit, offset)
		} else {
			// TODO: load user public thread messages by name
		}
	}

	if (is.notEmpty(ctx.state.messages)) {
		if (ctx.state.messages.error.error) {
			console.error(ctx.state.messages.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.messages = models.threadMessages(
			ctx.state.messages.messages,
			ctx.state.messages.paging,
			ctx.state.message,
			ctx.state.thread,
		)
	} else {
		ctx.state.messages = undefined
	}

	await next()

	console.log("<-- loadThreadMessages")
}

export default loadThreadMessages;
