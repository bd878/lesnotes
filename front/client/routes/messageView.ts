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
import SearchFormBuilder from '../builders/searchFormBuilder';

async function messageView(ctx) {
	ctx.log.info("--> messageView")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (!ctx.userAgent.isMobile) {
		content
			.addMessagesTree(
				tree
					.addThreadPath(ctx.state.threadPath || [] /* for root */)
					.addList(ctx.state.tree, 0)
			)
	}

	if (ctx.state.messageFeatures.hasFeatures()) {
		content.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))

		if (ctx.userAgent.isMobile) {
			layout.addNoFixHeight()
		}
	}

	layout
		.addTitle(ctx.state.message.title || ctx.state.message.text)
		.addFooter()
		.addContent(
			content
				.addMessageView(
					view
						.addRedirectUrl("")
						.addMessageEditHref(ctx.state.messageEditHref)
						.addMessage(ctx.state.message)
				)
				.addMessageHeader(
					messageHeader
						.addMessagePath(ctx.state.messagePath)
				)
				.addHeader(
					header
						.addNewNote(ctx.state.newNoteHref)
						.addSearch(searchForm)
				)
				.addControlPanel(panel.addAuth(auth.addLogout()))
		)

	ctx.body = layout.build()
	ctx.status = 200

	ctx.log.info("<-- messageView")
}

export default messageView;
