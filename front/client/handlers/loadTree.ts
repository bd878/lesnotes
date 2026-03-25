import api from '../api';
import * as is from '../third_party/is';

async function loadTree(ctx, next) {
	console.log("--> loadTree")

	ctx.state.tree = await api.readTreeJson(ctx.state.token, ctx.state.messageID, ctx.state.cwd.id, ctx.state.cwd.limit, ctx.state.cwd.offset, ctx.state.leaves)

	if (is.notEmpty(ctx.state.tree)) {
		if (ctx.state.tree.error.error) {
			console.error(ctx.state.tree.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}
		ctx.state.tree = ctx.state.tree.list
	} else {
		ctx.state.tree = {}
	}

	await next()

	console.log("<-- loadTree")
}

export default loadTree