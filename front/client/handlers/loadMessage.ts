import api from '../api';

async function loadMessage(ctx, next) {
	const id = parseInt(ctx.query.id) || 0

	if (id != 0)
		ctx.state.message = await api.readMessageJson(ctx.state.token, 0 /* me */, id)

	await next()
}

export default loadMessage
