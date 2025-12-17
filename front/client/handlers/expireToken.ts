
async function expireToken(ctx, next) {
	console.log("--> expireToken")

	ctx.set({"Set-Cookie": "token=\"\"; Expires=0; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})
	await next()

	console.log("<-- expireToken")
}

export default expireToken;
