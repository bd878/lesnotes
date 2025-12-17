import api from '../api';
import * as is from '../third_party/is';

async function loadSearchPath(ctx, next) {
	const token  = ctx.state.token
	const search = ctx.state.search

	console.log("--> loadSearchPath")

	if (is.notEmpty(token)) {
		if (is.notEmpty(search)) {
			if (!search.error.error) {
				ctx.state.searchPath = await api.searchMessagesPathJson(token, search.messages)
			} else {
				console.log("[loadSearchPath]: search error")
			}
		} else {
			console.log("[loadSearchPath]: search is empty")
		}

		console.log("[loadSearchPath]:", "searchPath", ctx.state.searchPath)
	}

	await next()

	console.log("<-- loadSearchPath")
}

export default loadSearchPath
