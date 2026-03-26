import * as is from '../third_party/is';
import PublicTranslationBuilder from '../builders/publicTranslationBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';

async function publicTranslation(ctx) {
	console.log("--> publicTranslation")

	const content = new PublicTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.msg == "comments") {
		content.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		content.addFilesView(ctx.state.message.files)
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			content.addFilesView(ctx.state.message.files)
		} else {
			content.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	content.addMessageNavigation(messageNavigation)
	content.addTranslations(ctx.state.messageName, ctx.state.message.translations)
	content.addTranslationView(ctx.state.translation)

	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- publicTranslation")
}

export default publicTranslation;
