import * as is from '../third_party/is';
import TranslationsBuilder from '../builders/translationsBuilder';

async function messageTranslations(ctx, next) {
	console.log("--> messageTranslations")

	const translations = new TranslationsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (is.notEmpty(ctx.state.messageFeatures) && is.notEmpty(ctx.state.trans)) {
		if (ctx.state.trans.mode == "new") {
			translations.addNewTranslation(ctx.state.messageID)
		} else if (ctx.state.trans.mode == "edit") {
			translations.addTranslationEditForm(ctx.state.messageID, ctx.state.translation)
		} else if (ctx.state.trans.mode == "view") {
			translations.addTranslationView(ctx.state.messageID, ctx.state.translation)
		} else {
			translations.addTranslationsList(ctx.state.messageID, ctx.state.message.translations)
		}

		ctx.state.messageFeatures.addTranslations(translations)
	}

	await next()

	console.log("<-- messageTranslations")
}

export default messageTranslations
