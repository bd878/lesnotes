import * as is from '../third_party/is';
import MessageEditViewBuilder from '../builders/messageEditViewBuilder';
import HomeBuilder from '../builders/homeBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function messageEdit(ctx) {
	console.log("--> messageEdit")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const messageForm = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)

	header.addNewNote()
	tree.addThreadPath(ctx.state.cwdPath)
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)
	messageForm
		.addFilesList(ctx.state.message.files)
		.addMessage(ctx.state.message)

	content.addMessagesTree(tree)
	content.addMessageEditForm(messageForm)
	content.addMessageFeatures(ctx.state.messageFeatures)
	auth.addLogout()
	content.addAuth(auth)
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
