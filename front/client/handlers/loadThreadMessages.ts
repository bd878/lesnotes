import api from '../api';
import * as is from '../third_party/is';

const limit = parseInt(LIMIT)

async function loadThreadMessages(ctx, next) {
	const userID = parseInt(ctx.params.user) || 0
	const id = is.notEmpty(ctx.state.thread) ? ctx.state.thread.ID : (parseInt(ctx.params.id) || 0);
	const token = ctx.state.token
	const params = new URLSearchParams(ctx.search)
	const offset = parseInt(params.get("offset")) || 0

	console.log("--> loadThreadMessages")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.messages = await api.readMessagesJson(token, id, 0, limit, offset)
		} else {
			// TODO: load private thread messages by name
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			ctx.state.messages = await api.readMessagesJson(token, id, userID, limit, offset)
		} else {
			// TODO: load user public thread messages by name
		}
	} else if (is.notEmpty(id) /* cannot read root (id=0) thread */) {
		// load public thread messages
		ctx.state.messages = await api.readMessagesJson("", id, 0, limit, offset)
	}

	if (is.notEmpty(ctx.state.messages)) {
		if (ctx.state.messages.error.error) {
			console.error(ctx.state.messages.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}
		ctx.state.messages = ctx.state.messages.messages
	} else {
		ctx.state.messages = undefined
	}

	await next()

	console.log("<-- loadThreadMessages")
}

export default loadThreadMessages;
