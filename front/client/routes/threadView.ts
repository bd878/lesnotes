import ThreadViewBuilder from '../builders/threadViewBuilder';
import HomeBuilder from '../builders/homeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';

async function threadView(ctx) {
	console.log("--> threadView")

	const thread = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	tree.addThreadPath(ctx.state.cwdPath)
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)

	header.addNewNote()

	thread.addThread(ctx.state.thread)

	content.addMessagesTree(tree)
	content.addThreadView(thread)
	content.addMessageHeader(messageHeader)
	content.addHeader(header)
	auth.addLogout()
	content.addAuth(auth)
	content.addControlPanel()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- threadView")
}

export default threadView;
