import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import PublicThreadMessageBuilder from '../builders/publicThreadMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new PublicThreadMessageBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout
		.addFooter()
		.addContent(
			content
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addHeader(header)
				.addMessageView(
					view
						.addDeleteRedirectUrl("/t/" + ctx.state.threadName + ctx.search)
						.addMessage(ctx.state.message)
				)
				.addMessagesTree(
					tree
						.addList(ctx.state.tree)
						.addThread(ctx.state.thread)
				)
				.addMessageFeatures(
					ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)
				)
		)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

export default publicThreadMessage
