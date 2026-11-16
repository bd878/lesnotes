
async function expireToken(ctx, next) {
	ctx.log.info("--> expireToken")

	ctx.set({"Set-Cookie": "token=; Expires=0; Max-Age=0; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})
	ctx.set({"Clear-Site-Data": "cookies"})
	await next()

	ctx.log.info("<-- expireToken")
}

export default expireToken;
