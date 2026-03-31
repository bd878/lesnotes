import NewMessageBuilder from '../builders/newMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function newMessage(ctx) {
	console.log("--> newMessage")

	const content = new NewMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	header.addNewNote()
	tree.addThreadPath(ctx.state.cwdPath)
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)

	content.addMessagesTree(tree)
	content.addFilesSelector([])
	content.addNewMessageForm(ctx.state.cwd.id)
	content.addMessageHeader(messageHeader)
	auth.addLogout()
	content.addAuth(auth)
	content.addControlPanel()
	content.addHeader(header)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- newMessage")
}

export default newMessage;
