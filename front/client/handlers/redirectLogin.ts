async function redirectLogin(ctx, next) {
	console.log("--> redirect login")
	ctx.redirect("/login" + ctx.search)
	ctx.status = 302
	console.log("<-- redirect login")
}

export default redirectLogin
