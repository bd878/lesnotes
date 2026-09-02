async function redirectHome(ctx, next) {
	ctx.log.info("--> redirectHome")
	ctx.redirect(ctx.router.url("home", {}, {query: ctx.query}))
	ctx.status = 302
	ctx.log.info("<-- redirectHome")
}

export default redirectHome
