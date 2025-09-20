import api from '../api';

async function getMe(ctx, next) {
	const resp = await api.authJson(ctx.state.token)
	if (resp.error.error || resp.expired) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	ctx.state.me = await api.getMeJson(ctx.state.token)
	if (ctx.state.me.error.error) {
		ctx.redirect("/login")
		ctx.status = 302
		return
	}

	await next()
}

export default getMe
