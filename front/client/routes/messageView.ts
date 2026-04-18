import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import HomeBuilder from '../builders/homeBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import TranslationsBuilder from '../builders/translationsBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
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
				.addMessageView(
					view
						.addRedirectUrl("")
						.addMessage(ctx.state.message)
				)
				.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
				.addMessageHeader(messageHeader.addMessagePath(ctx.state.messagePath))
				.addHeader(header.addNewNote("/home" + ctx.search))
				.addControlPanel(panel.addAuth(auth.addLogout()))
		)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- messageView")
}

export default messageView;
