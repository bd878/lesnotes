import ControlPanelBuilder from '../builders/controlPanelBuilder';
import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageViewBuilder from '../builders/messageViewBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';
import ThreadPathBuilder from '../builders/threadPathBuilder'
import MessagesListBuilder from '../builders/messagesListBuilder';

async function publicThread(ctx) {
	ctx.log.info("--> publicThread")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new PublicThreadBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const threadPath = new ThreadPathBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messagesList = new MessagesListBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.messageFeatures.hasFeatures()) {
		content.addThreadFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
	}

	layout
		.addTitle(ctx.state.message.title || ctx.state.message.text)
		.addFooter()
		.addContent(
			content
				.addMessagesTree(
					tree
						.addList(messagesList.addTree(ctx.state.tree))
						.addThreadPath(threadPath.addPath(ctx.state.threadPath))
				)
				.addMessageView(
					view
						.addMessageEditHref(ctx.state.messageEditHref)
						.addDeleteRedirectUrl(ctx.state.deleteRedirectUrl)
						.addMessage(ctx.state.message)
				)
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addHeader(ctx.state.isAuthed
					? header
						.addNewNote(ctx.state.newNoteHref)
						.addSearch(searchForm)
					: header.addSearch(searchForm.addThreadName(ctx.state.parentName))
				)
		)

	ctx.body = layout.build()
	ctx.status = 200

	ctx.log.info("<-- publicThread")
}

export default publicThread
