import * as is from '../third_party/is';
import MessageFeaturesBuilder from '../builders/messageFeaturesBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';

async function messageFeatures(ctx, next) {
	console.log("--> messageFeatures")

	const messageFeatures = new MessageFeaturesBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.nav == "comments") {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			messageNavigation.addAttachments()
		}

		messageNavigation.addTranslations()
		messageFeatures.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.nav == "files") {
		messageNavigation.addTranslations()
		messageNavigation.addAttachments()

		messageFeatures.addFilesView(ctx.state.message.files)
	} else if (ctx.state.nav == "trans") {
		messageNavigation.addComments()
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			messageNavigation.addAttachments()
		}
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			messageNavigation.addAttachments()
			messageNavigation.addComments()
			messageNavigation.addTranslations()
			messageFeatures.addFilesView(ctx.state.message.files)
		} else {
			messageNavigation.addTranslations()
			messageFeatures.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	messageFeatures.addNavigation(messageNavigation)

	ctx.state.messageFeatures = messageFeatures

	await next()

	console.log("<-- messageFeatures")
}

export default messageFeatures