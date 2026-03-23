import ThreadViewBuilder from '../builders/threadViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';

async function threadView(ctx) {
	console.log("--> threadView")

	const content = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	content.addNavigation()
	content.addControlPanel()
	content.addThreadView(ctx.state.thread)
	content.addMessagesTree(ctx.state.stack)
	content.addLogout()

	header.addSearch()

	layout.addSettings()
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- threadView")
}

export default threadView;
