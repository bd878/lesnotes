async function getToken(ctx, next) {
	const token = ctx.cookies.get("token")
	ctx.state.token = token

	console.log("--> getToken")

	await next()

	console.log("<-- getToken")
}

export default getToken
