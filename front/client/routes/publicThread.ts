import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const content = new PublicThreadBuilder(ctx.state.isAuthed, ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	tree.addList(ctx.state.tree)

	if (ctx.state.isAuthed) {
		auth.addLogout()
	} else {
		auth.addLogin()
	}

	header.addAuth(auth)
	content.addMessagesTree(tree)
	content.addHeader(header)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
