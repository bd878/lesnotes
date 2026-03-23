import NewTranslationBuilder from '../builders/newTranslationBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';

async function newTranslation(ctx) {
	console.log("--> newTranslation")

	const content = new NewTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	content.addMessagesTree(ctx.state.stack)
	content.addNewTranslationForm(ctx.state.messageID)
	content.addLogout()

	header.addSearch()

	layout.addSettings()
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- newTranslation")
}

export default newTranslation;
