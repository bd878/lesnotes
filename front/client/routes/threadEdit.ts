import ThreadEditBuilder from '../builders/threadEditBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import AuthBuilder from '../builders/authBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';

async function threadEdit(ctx) {
	console.log("--> threadEdit")

	const content = new ThreadEditBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	tree.addThreadPath(ctx.state.cwdPath)
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)

	auth.addLogout()
	header.addNewNote()

	content.addMessagesTree(tree)
	content.addAuth(auth)
	content.addThreadEditForm(ctx.state.thread)
	content.addMessageHeader(messageHeader)
	content.addHeader(header)
	content.addControlPanel()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- threadEdit")
}

export default threadEdit
