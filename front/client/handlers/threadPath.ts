import * as is from '../third_party/is';
import {EmptyMessage} from '../api/models/message';
import readPathJson from '../api/readPathJson';

async function threadPath(ctx, next) {
	const token = ctx.state.token

	ctx.log.info("--> threadPath")

	if (is.notEmpty(ctx.state.parentName)) {
		const result = await readPathJson(token, 0, ctx.state.parentName)

		if (is.notEmpty(result)) {
			if (result.error.error) {
				ctx.log.error(result.error)
				ctx.body = "error"
				ctx.status = 400;
				return
			}

			if (ctx.state.isAuthed) {
				result.path.push(EmptyMessage /* root */)
			}

			result.path.reverse()

			ctx.state.threadPath = result.path
		} else {
			ctx.state.threadPath = []
		}
	} else  if (is.notEmpty(ctx.state.message)) {
		ctx.state.threadPath = [ctx.state.isAuthed && EmptyMessage, ctx.state.message].filter(Boolean)
	}

	await next()

	ctx.log.info("<-- threadPath")
}

export default threadPath