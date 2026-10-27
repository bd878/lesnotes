import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import MessageViewBuilder from '../builders/messageViewBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import ThreadPathBuilder from '../builders/threadPathBuilder';
import SearchBuilder from '../builders/searchBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';

async function searchMessageView(ctx) {
	ctx.log.info("--> searchMessageView")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const threadPath = new ThreadPathBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const content = new SearchBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.messageFeatures.hasFeatures()) {
		content.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
	}

	layout
		.addTitle(ctx.state.searchQuery)
		.addFooter()
		.addContent(
			content
				.addThreadPath(threadPath.addPath(ctx.state.threadPath || []))
				.addHeader(
					ctx.state.isAuthed
						? header.addSearch(searchForm).addNewNote(ctx.state.newNoteHref)
						: header.addSearch(searchForm)
				)
				.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
				.addMessagesList(ctx.state.list)
				.addMessageView(
					view
						.addDeleteRedirectUrl(ctx.state.deleteRedirectUrl)
						.addMessageEditHref(ctx.state.messageEditHref)
						.addMessage(ctx.state.message)
				)
		)

	ctx.body = layout.build()
	ctx.status = 200;

	ctx.log.info("<-- searchMessageView")
}

export default searchMessageView
