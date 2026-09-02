import * as is from '../third_party/is';
import api from '../api';

async function authed(ctx, next) {
	ctx.log.info("--> authed")

	if (is.empty(ctx.state.token)) {
		ctx.redirect('/login' + ctx.search)
		ctx.status = 302
	} else {
		const resp = await api.authJson(ctx.state.token)
		if (resp.error.error || resp.expired) {
			ctx.log.error(resp.error)
			ctx.redirect('/login' + ctx.search)
			ctx.status = 302
		} else {
			await next()
		}
	}

	ctx.log.info("<-- authed")
}

export default authed
