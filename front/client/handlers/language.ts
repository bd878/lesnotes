async function loadLanguage(ctx, next) {
	// may be mutated by getMe
	ctx.state.lang = ctx.acceptsLanguages(["ru", "en", "fr", "de"])
	if (ctx.state.lang == "")
		ctx.state.lang = "en"
	await next()
}

export default loadLanguage
