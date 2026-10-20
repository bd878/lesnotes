import searchMessagesJson from '../api/searchMessagesJson'
import * as is from '../third_party/is';

async function searchMessages(ctx, next) {
	const token = ctx.state.token
	const isPublic = ctx.state.isAuthed ? -1 : 1;

	ctx.log.info("--> searchMessages")

	const resp = await searchMessagesJson(token, ctx.state.search.query, parseInt(ctx.state.search.id) || 0, ctx.state.search.name, isPublic)
	if (is.notEmpty(resp) && resp.error.error) {
		ctx.log.error(resp.error)
		ctx.body = "error"
		ctx.status = 400
		return
	}

	ctx.state.list = resp.list

	await next()

	ctx.log.info("<-- searchMessages")
}

export default searchMessages