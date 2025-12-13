async function redirectHome(ctx, next) {
	console.log("--> redirect home")
	ctx.redirect("/home" + ctx.search)
	ctx.status = 302
	console.log("<-- redirect home")
}

export default redirectHome
