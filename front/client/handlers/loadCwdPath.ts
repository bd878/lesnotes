import * as is from '../third_party/is';
import {EmptyMessage} from '../api/models/message';
import readPathJson from '../api/readPathJson';

async function loadCwdPath(ctx, next) {
	const token = ctx.state.token

	ctx.log.info("--> loadCwdPath")

	if (is.notEmpty(ctx.state.cwd) && (ctx.state.cwd.id != 0) /* not root */) {
		const result = await readPathJson(token, ctx.state.cwd.id, "")

		if (is.notEmpty(result)) {
			if (result.error.error) {
				ctx.log.error(result.error)
				ctx.body = "error"
				ctx.status = 400;
				return
			}

			result.path.push(EmptyMessage /* root */)

			result.path.reverse()

			ctx.state.cwdPath = result.path
		} else {
			ctx.state.cwdPath = []
		}
	} else {
		ctx.state.cwdPath = []
	}

	await next()

	ctx.log.info("<-- loadCwdPath")
}

export default loadCwdPath
