import * as is from '../third_party/is';
import api from '../api';

async function authed(ctx, next) {
	console.log("--> authed")

	if (is.empty(ctx.state.token)) {
		ctx.redirect('/login' + ctx.search)
		ctx.status = 302
	} else {
		const resp = await api.authJson(ctx.state.token)
		if (resp.error.error || resp.expired) {
			ctx.redirect('/login' + ctx.search)
			ctx.status = 302
		} else {
			await next()
		}
	}

	console.log("<-- authed")
}

export default authed
