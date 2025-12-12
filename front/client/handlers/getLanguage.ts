async function getLanguage(ctx, next) {
	console.log("--> getLanguage")

	switch (ctx.query.lang) {
	case "ru":
	case "en":
	case "fr":
	case "de":
		ctx.state.lang = ctx.query.lang;
		break;
	default:
		ctx.state.lang = ctx.acceptsLanguages(["ru", "en", "fr", "de"])
		if (!ctx.state.lang) {
			ctx.state.lang = "en"
		}
	}

	await next()
	console.log("<-- getLanguage")
}

export default getLanguage
