async function getToken(ctx, next) {
	const token = ctx.cookies.get("token")
	ctx.state.token = token
	await next()
}

export default getToken
