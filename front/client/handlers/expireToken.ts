
async function expireToken(ctx, next) {
	ctx.log.info("--> expireToken")

	ctx.set({"Set-Cookie": "token=\"\"; Expires=0; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})
	await next()

	ctx.log.info("<-- expireToken")
}

export default expireToken;
