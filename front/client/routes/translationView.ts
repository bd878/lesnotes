import * as is from '../third_party/is';
import TranslationViewBuilder from '../builders/translationViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function translationView(ctx) {
	console.log("--> translationView")

	const content = new TranslationViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.msg == "comments") {
		content.addMessageNavigation()
		content.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		content.addMessageNavigation()
		content.addFilesView(ctx.state.message.files)
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			content.addMessageNavigation()
			content.addFilesView(ctx.state.message.files)
		} else {
			content.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	content.addNavigation()
	content.addControlPanel()
	content.addMessagesTree(ctx.state.stack)
	content.addNewTranslation(ctx.state.message.ID)
	content.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	content.addTranslationView(ctx.state.message.ID, ctx.state.translation)
	content.addLogout()

	header.addSearch()

	layout.addSettings(settings)
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- translationView")
}

export default translationView;
