import ThreadViewBuilder from '../builders/threadViewBuilder';
import HomeBuilder from '../builders/homeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';

async function threadView(ctx) {
	console.log("--> threadView")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const thread = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout
		.addFooter()
		.addContent(			
			content
				.addMessagesTree(
					tree
						.addThreadPath(ctx.state.cwdPath)
						.addList(ctx.state.tree)
				)
				.addThreadView(thread.addThread(ctx.state.thread))
				.addMessageHeader(messageHeader.addMessagePath(ctx.state.messagePath))
				.addHeader(header.addNewNote("/home" + ctx.search))
				.addControlPanel(panel.addAuth(auth.addLogout()))
		)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- threadView")
}

export default threadView;
