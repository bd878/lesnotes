async function language(ctx, next) {
	// may be mutated by getMe
	console.log("--> language")
	ctx.state.lang = ctx.acceptsLanguages(["ru", "en", "fr", "de"])
	if (!ctx.state.lang) {
		ctx.state.lang = "en"
	}
	await next()
	console.log("<-- language")
}

export default language
