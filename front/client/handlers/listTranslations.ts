import listTranslationsJson, { EmptyListTranslations } from '../api/listTranslationsJson'
import * as is from '../third_party/is';

async function listTranslations(ctx, next) {
	const id = parseInt(ctx.query.id) || parseInt(ctx.params.id) || 0
	const name = ctx.params.messageName || ""
	const token = ctx.state.token

	console.log("--> listTranslations")

	ctx.state.translations = await listTranslationsJson(token, id, name)
	if (is.notEmpty(ctx.state.translations)) {
		if (ctx.state.translations.error.error) {
			console.error(ctx.state.translations.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.translations = ctx.state.translations.translations
	} else {
		ctx.state.translations = EmptyListTranslations.translations
	}

	await next()

	console.log("<-- listTranslations")
}

export default listTranslations
