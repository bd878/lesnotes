import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import PublicThreadMessageBuilder from '../builders/publicThreadMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';
import ThreadPathBuilder from '../builders/threadPathBuilder'
import MessagesListBuilder from '../builders/messagesListBuilder';

async function publicThreadMessage(ctx) {
	ctx.log.info("--> publicThreadMessage")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new PublicThreadMessageBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const threadPath = new ThreadPathBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messagesList = new MessagesListBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.messageFeatures.hasFeatures()) {
		content.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
	}

	layout
		.addTitle(ctx.state.message.title || ctx.state.message.text)
		.addFooter()
		.addContent(
			content
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addHeader(ctx.state.isAuthed
					? header.addNewNote(ctx.state.newNoteHref).addSearch(searchForm)
					: header.addSearch(searchForm.addThreadName(ctx.state.parentName))
				)
				.addMessageView(
					view
						.addMessageEditHref(ctx.state.messageEditHref)
						.addDeleteRedirectUrl(ctx.state.deleteRedirectUrl)
						.addMessage(ctx.state.message)
				)
				.addMessagesTree(
					tree
						.addList(messagesList.addTree(ctx.state.tree))
						.addThreadPath(threadPath.addPath(ctx.state.threadPath))
				)
		)

	ctx.body = layout.build()
	ctx.status = 200

	ctx.log.info("<-- publicThreadMessage")
}

export default publicThreadMessage
