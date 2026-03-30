import * as is from '../third_party/is';
import {EmptyMessage} from '../api/models/message';
import readPathJson from '../api/readPathJson';

async function loadCwdPath(ctx, next) {
	const token = ctx.state.token

	console.log("--> loadCwdPath")

	console.log("=== CWD ===\n", ctx.state.cwd)

	if (is.notEmpty(ctx.state.cwd) && ctx.state.cwd.id != 0 /* not root */) {
		const result = await readPathJson(token, ctx.state.cwd.id)

		if (is.notEmpty(result)) {
			if (result.error.error) {
				console.error(result.error)
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

	console.log("=== PATH===\n", ctx.state.cwdPath)

	await next()

	console.log("<-- loadCwdPath")
}

export default loadCwdPath
