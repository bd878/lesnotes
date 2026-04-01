import * as is from '../third_party/is';
import MessageFeaturesBuilder from '../builders/messageFeaturesBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';
import TranslationsBuilder from '../builders/translationsBuilder';

async function messageFeatures(ctx, next) {
	console.log("--> messageFeatures")

	const messageFeatures = new MessageFeaturesBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const translations = new TranslationsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
		messageNavigation.addAttachments()
	}

	messageNavigation.addComments()

	if (ctx.state.isAuthed || ctx.state.message.translations.length > 0) {
		messageNavigation.addTranslations()
	}

	if (ctx.state.nav == "comments") {
		messageFeatures.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.nav == "files") {

		messageFeatures.addFilesView(ctx.state.message.files)

	} else if (ctx.state.nav == "trans") {

		if (is.notEmpty(ctx.state.trans)) {
			if (ctx.state.trans.mode == "new") {
				translations.addTranslationNewForm(ctx.state.messageID)
			} else if (ctx.state.trans.mode == "edit") {
				translations.addTranslationEditForm(ctx.state.messageID, ctx.state.translation)
			} else if (ctx.state.trans.mode == "view") {
				translations.addTranslationView(ctx.state.messageID, ctx.state.translation)
			} else {
				translations.addTranslationsList(ctx.state.messageID, ctx.state.message.translations)
				translations.addNewTranslation(ctx.state.messageID)
			}

			messageFeatures.addTranslations(translations)
		}

	} else {
		messageFeatures.addComments(ctx.state.message.ID, ctx.state.comments)
	}

	ctx.state.messageFeatures = messageFeatures
	ctx.state.messageNavigation = messageNavigation

	await next()

	console.log("<-- messageFeatures")
}

export default messageFeatures