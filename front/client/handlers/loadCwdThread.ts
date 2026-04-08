import api from '../api';
import * as is from '../third_party/is';

async function loadCwdThread(ctx, next) {
	const token = ctx.state.token

	console.log("--> loadCwdThread")

	if (is.notEmpty(ctx.state.cwd.name) && is.empty(ctx.state.cwd.id)) {
		const resp = await api.readThreadJson(token, 0, 0, ctx.state.cwd.name)
		if (is.notEmpty(resp) && resp.error.error) {
			console.error(resp.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}
		ctx.state.cwd.id = resp.thread.ID
	}

	await next()

	console.log("<-- loadCwdThread")
}

export default loadCwdThread
