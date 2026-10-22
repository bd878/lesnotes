import NewMessageBuilder from '../builders/newMessageBuilder'
import HomeBuilder from '../builders/homeBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';

async function newRootMessage(ctx) {
	ctx.log.info("--> newMessage")

	if (!ctx.userAgent.isMobile) {
		ctx.redirect(ctx.router.url("home", {}, {query: ctx.query}))
		return
	}

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const messageForm = new NewMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.tree.count > 0) {
		header.addSearch(searchForm)
	}

	layout
		.addFooter()
		.addContent(
			content
				.addNewMessageForm(
					messageForm.addFilesList([])
				)
				.addControlPanel(
					panel.addAuth(auth.addLogout())
				)
				.addHeader(header.addNewNote(ctx.state.newNoteHref))
		)

	ctx.body = layout.build()
	ctx.status = 200;

	ctx.log.info("<-- newMessage")
}

export default newRootMessage
