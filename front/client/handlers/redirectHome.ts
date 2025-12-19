async function redirectHome(ctx, next) {
	console.log("--> redirectHome")
	ctx.redirect(ctx.router.url("home", {}, {query: ctx.query}))
	ctx.status = 302
	console.log("<-- redirectHome")
}

export default redirectHome
