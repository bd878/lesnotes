
async function expireToken(ctx, next) {
	console.log("--> expire token")

	ctx.set({"Set-Cookie": "token=\"\"; Expires=0; HttpOnly; Path=/; Secure; Domain=" + `${DOMAIN}`})
	await next()

	console.log("<-- expire token")
}

export default expireToken;
