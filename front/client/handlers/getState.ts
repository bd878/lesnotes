async function getState(ctx, next) {
	console.log("--> getState")

	ctx.state.fontSize = getFontSize(ctx)
	ctx.state.lang = getLanguage(ctx)
	ctx.state.theme = getTheme(ctx)
	ctx.state.thread = getThread(ctx)

	await next()

	console.log("<-- getState")
}

function getFontSize(ctx) {
	switch (ctx.query.size) {
	case "small":
	case "medium":
	case "large":
		return ctx.query.size
	default:
		return "medium"
	}
}

function getLanguage(ctx) {
	switch (ctx.query.lang) {
	case "ru":
	case "en":
	case "fr":
	case "de":
		return ctx.query.lang;
	default:
		return ctx.acceptsLanguages(["ru", "en", "fr", "de"]) || "en"
	}
}

function getTheme(ctx) {
	switch (ctx.query.theme) {
	case "dark":
		return "dark"
	case "light":
		return "light"
	default:
		return "light"
	}
}

function getThread(ctx) {
	return parseInt(ctx.query.cwd) || 0
}

export default getState;
