import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import LogoutBuilder from '../builders/logoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const logout = new LogoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.msg == "comments") {
		messageNavigation.addTranslations()
		content.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		messageNavigation.addAttachments()
		content.addFilesView(ctx.state.message.files)
	} else if (ctx.state.msg == "translations") {
		messageNavigation.addComments()
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			messageNavigation.addAttachments()
			messageNavigation.addComments()
			content.addFilesView(ctx.state.message.files)
		}
		messageNavigation.addTranslations()

		content.addComments(ctx.state.message.ID, ctx.state.comments)
	}

	header.addNewNote()
	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)
	messageHeader.addThreadLink(ctx.state.message.ID)

	content.addMessagesTree(tree)
	content.addMessageNavigation(messageNavigation)
	content.addMessageView(ctx.state.me.ID, ctx.state.message)
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
