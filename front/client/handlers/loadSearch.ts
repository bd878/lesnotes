import api from '../api';
import redirectHome from './redirectHome';

async function loadSearch(ctx, next) {
	console.log("--> loadSearch")

	const response = await api.searchMessagesJson(ctx.state.token, ctx.state.query)
	if (response.error.error) {
		console.log(response.error)
		ctx.body = "error"
		ctx.state = 400
		// TODO: show error screen
		return await redirectHome(ctx, next)
	} else {
		ctx.state.messages = response.messages
	}

	await next()

	console.log("<-- loadSearch")
}

export default loadSearch
