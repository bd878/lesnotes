import api from '../api';
import * as is from '../third_party/is';

async function getMe(ctx, next) {
	console.log("--> getMe")

	ctx.state.me = await api.getMeJson(ctx.state.token)
	if (ctx.state.me.error.error) {
		ctx.redirect("/login")
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

	console.log("<-- getMe")
}

export default getMe
