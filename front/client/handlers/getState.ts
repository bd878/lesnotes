import type {IDLimitOffset} from '../types'

const defaultLimit = parseInt(LIMIT) || 10

async function getState(ctx, next) {
	console.log("--> getState")

	ctx.state.fontSize = getFontSize(ctx)
	ctx.state.lang = getLanguage(ctx)
	ctx.state.theme = getTheme(ctx)
	ctx.state.thread = getThread(ctx)
	ctx.state.nav = getMessageView(ctx)
	ctx.state.trans = getTranslation(ctx)
	ctx.state.cwd = getCwd(ctx)
	ctx.state.messageID = getMessageID(ctx)
	ctx.state.messageName = ctx.params.messageName || ""
	ctx.state.threadID = getThreadID(ctx)
	ctx.state.leaves = getLeaves(ctx)
	ctx.state.token = getToken(ctx)

	await next()

	console.log("<-- getState")
}

function getToken(ctx): string {
	return ctx.cookies.get("token")
}

function getMessageID(ctx): number {
	return parseInt(ctx.params.id) || 0
}

function getThreadID(ctx): number {
	return parseInt(ctx.params.id) || 0
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

function getTranslation(ctx) {
	const [lang = "", mode = ""] = ((new URLSearchParams(ctx.request.search)).get("trans") || "").split(",")
	const result = { lang: "", mode: "" }

	switch (lang) {
	case "ru":
	case "en":
	case "fr":
	case "de":
		result.lang = lang
		break
	case "new":
		result.mode = "new"
		break;
	}

	switch (mode) {
	case "edit":
	case "view":
		result.mode = mode
		break
	}

	return result
}

function getMessageView(ctx) {
	switch (ctx.query.nav) {
	case "files":
	case "comments":
	case "trans":
	// TODO: case "trans":
		return ctx.query.nav
	default:
		return ""
	}
}

function idLimitOffset(numbers: number[]): IDLimitOffset {
	const [id = 0, limit = defaultLimit, offset = 0] = numbers
	return {id, limit, offset}
}

function getCwd(ctx) {
	const params = new URLSearchParams(ctx.request.search)

	return idLimitOffset([parseInt(params.get("cwd")) || 0, ...(params.get(params.get("cwd") || "0") || "").split(",").map(parseFloat).filter(v => !isNaN(v))])
}

function getLeaves(ctx): IDLimitOffset[] {
	const result = []

	for (const [key, value] of new URLSearchParams(ctx.request.search)) {
		const threadID = parseInt(key)
		if (!isNaN(threadID)) {
			result.push(idLimitOffset([threadID, ...(value || "").split(",").map(parseFloat).filter(v => !isNaN(v))]))
		}
	}

	return result
}

function getThread(ctx) {
	return parseInt(ctx.query.cwd) || 0
}

export default getState;
