import * as is from '../third_party/is';
import MessageEditViewBuilder from '../builders/messageEditViewBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import LogoutBuilder from '../builders/logoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function messageEdit(ctx) {
	console.log("--> messageEdit")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const logout = new LogoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	ctx.state.messageNavigation.addAttachments()
	ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)

	header.addNewNote()
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)
	messageHeader.addThreadLink(ctx.state.message.ID)

	content.addMessagesTree(tree)
	content.addMessageEditForm(ctx.state.message)
	content.addMessageFeatures(ctx.state.messageFeatures)
	content.addLogout(logout)
	content.addMessageHeader(messageHeader)
	content.addHeader(header)
	content.addControlPanel()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- messageEdit")
}

export default messageEdit;
