async function getFontSize(ctx, next) {
	// may be mutated by getMe
	console.log("--> getFontSize")
	switch (ctx.query.size) {
	case "small":
	case "medium":
	case "large":
		ctx.state.fontSize = ctx.query.size
		break;
	default:
		ctx.state.fontSize = "medium"
	}
	await next()
	console.log("<-- getFontSize")
}

export default getFontSize
