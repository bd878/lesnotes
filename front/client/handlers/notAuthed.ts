import api from '../api';

async function notAuthed(ctx, next) {
	console.log("--> not authed")

	const resp = await api.authJson(ctx.state.token)
	if (resp.error.error || resp.expired) {
		await next()
	} else {
		ctx.redirect('/home')
		ctx.status = 302
	}

	console.log("<-- not authed")
}

export default notAuthed
