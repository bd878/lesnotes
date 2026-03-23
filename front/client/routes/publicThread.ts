import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const content = new PublicThreadBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	content.addMessagesList(ctx.params.threadName /* TODO: use from load_path, ctx.thread.name is message name now, but thread name required */, ctx.state.messages)

	if (ctx.state.me.ID) {
		content.addLogout()
	} else {
		content.addSignup()
	}

	header.addSearch()

	layout.addSettings()
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
