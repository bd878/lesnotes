import api from '../api';

async function getMe(ctx, next) {
	console.log("--> getMe")

	ctx.state.me = await api.getMeJson(ctx.state.token)
	if (ctx.state.me.error.error) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	ctx.state.me = ctx.state.me.user

	if (ctx.state.me.lang) {
		ctx.state.lang = ctx.state.me.lang
	}

	console.log(ctx.state.me)

	await next()

	console.log("<-- getMe")
}

export default getMe
