import type { Message, File } from '../api/models'
import type { FileWithMime } from '../types'
import readTranslationJson, { EmptyReadTranslation } from '../api/readTranslationJson';
import * as is from '../third_party/is';

async function loadTranslation(ctx, next) {
	const id = parseInt(ctx.query.id) || parseInt(ctx.params.id) || 0
	const name = ctx.params.messageName || ""
	const lang = ctx.params.lang || ""
	const token = ctx.state.token

	console.log("--> loadTranslation")

	if (is.empty(lang)) {
		console.error("empty lang")
		ctx.body = "error"
		ctx.status = 400
		return
	}

	ctx.state.translation = await readTranslationJson(token, id, lang, "")
	if (is.notEmpty(ctx.state.translation)) {
		if (ctx.state.translation.error.error) {
			console.error(ctx.state.translation.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.translation = ctx.state.translation.translation
	} else {
		ctx.state.translation = EmptyReadTranslation.translation
	}

	await next()

	console.log("<-- loadTranslation")
}

export default loadTranslation
