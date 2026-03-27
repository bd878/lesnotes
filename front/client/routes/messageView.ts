import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import LogoutBuilder from '../builders/logoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import TranslationsBuilder from '../builders/translationsBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const logout = new LogoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)

	header.addNewNote()
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)
	if (is.notEmpty(ctx.state.thread) && is.notEmpty(ctx.state.thread.ID)) {
		messageHeader.addThreadLink(ctx.state.message.ID)
	} else {
		messageHeader.addEditThreadLink(ctx.state.message.ID)
	}

	content.addMessagesTree(tree)
	content.addMessageView(ctx.state.me.ID, ctx.state.message)
	content.addMessageFeatures(ctx.state.messageFeatures)
	content.addLogout(logout)
	content.addMessageHeader(messageHeader)
	content.addHeader(header)
	content.addControlPanel()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- messageView")
}

export default messageView;
