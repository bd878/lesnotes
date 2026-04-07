import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const content = new PublicThreadBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout
		.addFooter()
		.addContent(
			content
				.addMessagesTree(
					tree
						.addList(ctx.state.tree)
						.addThread(ctx.state.thread)
				)
				.addHeader(
					header.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin())
				)
		)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
