import TranslationEditViewBuilder from '../builders/translationEditViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function translationEdit(ctx) {
	console.log("--> translationEdit")

	const content = new TranslationEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	content.addNavigation()
	content.addControlPanel()
	content.addMessagesTree(ctx.state.stack)
	content.addFilesView(ctx.state.message.files)
	content.addNewTranslation(ctx.state.message.ID)
	content.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	content.addTranslationEditForm(ctx.state.message.ID, ctx.state.translation)
	content.addLogout()

	header.addSearch()

	layout.addSettings(settings)
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- translationEdit")
}

export default translationEdit;
