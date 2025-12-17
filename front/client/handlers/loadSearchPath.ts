import api from '../api';

async function loadSearchPath(ctx, next) {
	const token  = ctx.state.token

	console.log("--> loadSearchPath")

	const response = await api.searchMessagesPathJson(token, ctx.state.messages)
	if (response.error.error) {
		console.log(response.error)
		ctx.body = "error"
		ctx.state = 400
		return
	}

	ctx.state.messages = response.messages

	await next()

	console.log("<-- loadSearchPath")
}

export default loadSearchPath
