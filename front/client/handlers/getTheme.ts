async function getTheme(ctx, next) {
	console.log("--> getTheme")

	switch (ctx.query.theme) {
	case "dark":
		ctx.state.theme = "dark"
		break;
	case "light":
		ctx.state.theme = "light"
		break;
	default:
		ctx.state.theme = "light"
	}

	await next()

	console.log("<-- getTheme")
}

export default getTheme
