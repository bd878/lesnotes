import * as is from '../third_party/is';
import api from '../api';

async function notAuthed(ctx, next) {
	console.log("--> notAuthed")

	if (is.undef(ctx.state.token)) {
		await next()
	} else {
		const resp = await api.authJson(ctx.state.token)
		if (resp.error.error || resp.expired) {
			await next()
		} else {
			ctx.redirect('/home' + ctx.search)
			ctx.status = 302
		}
	}

	console.log("<-- notAuthed")
}

export default notAuthed
