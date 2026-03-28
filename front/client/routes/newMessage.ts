import NewMessageBuilder from '../builders/newMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import LogoutBuilder from '../builders/logoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function newMessage(ctx) {
	console.log("--> newMessage")

	const content = new NewMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const logout = new LogoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	header.addNewNote()
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)

	content.addMessagesTree(tree)
	content.addFilesSelector([])
	content.addNewMessageForm(ctx.state.messageID)
	content.addMessageHeader(messageHeader)
	content.addLogout(logout)
	content.addControlPanel()
	content.addHeader(header)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- newMessage")
}

export default newMessage;
