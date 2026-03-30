import ThreadViewBuilder from '../builders/threadViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import LogoutBuilder from '../builders/logoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';

async function threadView(ctx) {
	console.log("--> threadView")

	const content = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const logout = new LogoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	tree.addList(ctx.state.tree)

	messageHeader.addMessagePath(ctx.state.messagePath)

	header.addNewNote()

	content.addMessagesTree(tree)
	content.addThreadView(ctx.state.thread)
	content.addMessageHeader(messageHeader)
	content.addHeader(header)
	content.addLogout(logout)
	content.addControlPanel()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- threadView")
}

export default threadView;
