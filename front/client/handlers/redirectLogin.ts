async function redirectLogin(ctx, next) {
	console.log("--> redirectLogin")
	ctx.redirect("/login" + ctx.search)
	ctx.status = 302
	console.log("<-- redirectLogin")
}

export default redirectLogin
