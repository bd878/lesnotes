async function redirectLogin(ctx, next) {
	ctx.log.info("--> redirectLogin")
	ctx.redirect("/login" + ctx.search)
	ctx.status = 302
	ctx.log.info("<-- redirectLogin")
}

export default redirectLogin
