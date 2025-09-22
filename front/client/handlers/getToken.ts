async function getToken(ctx, next) {
	console.log("--> getToken")

	const token = ctx.cookies.get("token")
	ctx.state.token = token
	await next()

	console.log("<-- getToken")
}

export default getToken
