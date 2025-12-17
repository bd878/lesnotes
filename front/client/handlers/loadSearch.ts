import api from '../api';
import * as is from '../third_party/is';

async function loadSearch(ctx, next) {
	const query  = ctx.query.query || ""
	const token  = ctx.state.token

	console.log("--> loadSearch")

	if (is.notEmpty(token)) {
		ctx.state.search = await api.searchMessagesJson(token, query)
	}

	await next()

	console.log("<-- loadSearch")
}

export default loadSearch
