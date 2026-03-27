import type { Message, File } from '../api/models'
import type { FileWithMime } from '../types'
import readTranslationJson, { EmptyReadTranslation } from '../api/readTranslationJson';
import * as is from '../third_party/is';

async function loadTranslation(ctx, next) {   
	console.log("--> loadTranslation")

	if (is.notEmpty(ctx.state.trans) && is.notEmpty(ctx.state.trans.lang)) {
		ctx.state.translation = await readTranslationJson(ctx.state.token, ctx.state.messageID, ctx.state.trans.lang, ctx.state.messageName)
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
	}

	await next()

	console.log("<-- loadTranslation")
}

export default loadTranslation
