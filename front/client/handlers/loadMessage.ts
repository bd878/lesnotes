import api from '../api';
import * as is from '../third_party/is';

async function loadMessage(ctx, next) {
	const id = parseInt(ctx.query.id) || parseInt(ctx.params.id) || 0
	const name = ctx.params.name || ""
	const userID = parseInt(ctx.params.user) || 0
	const token = ctx.state.token

	console.log("--> loadMessage")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.message = await api.readMessageJson(token, 0 /* me */, id)
		} else if (is.notEmpty(name)) {
			ctx.state.message = await api.readMessageJson(token, 0 /* me */, 0, name /* public name */)
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			ctx.state.message = await api.readMessageJson("", userID, id)
		}
	} else if (is.notEmpty(name)) {
		ctx.state.message = await api.readMessageJson("", 0 /* public */, 0, name /* public name */)
	}

	if (is.notEmpty(ctx.state.message)) {
		if (ctx.state.message.error.error) {
			console.error(ctx.state.message.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}

		ctx.state.message = ctx.state.message.message
	} else {
		ctx.state.message = undefined
	}

	await next()

	console.log("<-- loadMessage")
}

export default loadMessage
