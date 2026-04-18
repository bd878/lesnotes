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
	const content = new PublicThreadMessageBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	const params = new URLSearchParams(ctx.search)
	params.set("cwd", ctx.state.message.ID)

	layout
		.addFooter()
		.addContent(
			content
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addHeader(ctx.state.isAuthed ? header.addNewNote("/home?" + params.toString()) : header)
				.addMessageView(
					view
						.addDeleteRedirectUrl("/" + ctx.state.parentName + ctx.search)
						.addMessage(ctx.state.message)
				)
				.addMessagesTree(
					tree
						.addList(ctx.state.tree)
						.addMessage(ctx.state.parentMessage)
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
