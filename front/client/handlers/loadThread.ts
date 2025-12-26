import api from '../api';
import * as is from '../third_party/is';

async function loadThread(ctx, next) {
	const token = ctx.state.token
	const userID = parseInt(ctx.params.user) || 0
	const id = parseInt(ctx.query.id) || parseInt(ctx.params.id) || 0
	const name = ctx.params.threadName || ""

	console.log("--> loadThread")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.thread = await api.readThreadJson(token, 0 /* me */, id)
		} else if (is.notEmpty(name)) {
			ctx.state.thread = await api.readThreadJson(token, 0 /* me */, 0, name /* public name */)
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			ctx.state.thread = await api.readThreadJson("", userID, id)
		}
	} else if (is.notEmpty(name)) {
		ctx.state.thread = await api.readThreadJson("", 0 /* public */, 0, name /* public name */)
	}

	if (is.notEmpty(ctx.state.thread)) {
		if (ctx.state.thread.error.error) {
			console.error(ctx.state.thread.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.thread = ctx.state.thread.thread
	} else {
		ctx.state.thread = undefined
	}

	await next()

	console.log("<-- loadThread")
}

export default loadThread
