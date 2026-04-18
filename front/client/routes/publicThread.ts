import ControlPanelBuilder from '../builders/controlPanelBuilder';
import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageViewBuilder from '../builders/messageViewBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new PublicThreadBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	const params = new URLSearchParams(ctx.search)
	params.set("cwd", ctx.state.message.ID)

	layout
		.addFooter()
		.addContent(
			content
				.addMessagesTree(
					tree
						.addList(ctx.state.tree)
						.addMessage(ctx.state.message)
				)
				.addMessageView(
					view
						.addDeleteRedirectUrl("/" + ctx.state.messageName + ctx.search)
						.addMessage(ctx.state.message)
				)
				.addThreadFeatures(
					ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)
				)
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addHeader(ctx.state.isAuthed ? header.addNewNote("/home?" + params.toString()) : header)
		)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
