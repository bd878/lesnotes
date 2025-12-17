import api from '../api';
import * as is from '../third_party/is';

function log(f) {
	return async function getMeLog(ctx, next) {
		console.log("--> getMe")
		await f(ctx, next)
		console.log("<-- getMe")
	}
}

async function getMe(ctx, next) {
	ctx.state.me = await api.getMeJson(ctx.state.token)
	if (ctx.state.me.error.error) {
		ctx.redirect("/login" + ctx.search)
		ctx.status = 302
		return
	}

	ctx.state.me = ctx.state.me.user

	if (is.empty(ctx.state.me)) {
		console.error("no me")
		ctx.status = 500
		return
	}

	await next()
}

export default log(getMe)
