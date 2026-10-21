import NewMessageBuilder from '../builders/newMessageBuilder'
import HomeBuilder from '../builders/homeBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';

async function newThreadMessage(ctx) {
	ctx.log.info("--> newThreadMessage")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const messageForm = new NewMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout
		.addFooter()
		.addContent(
			content
				.addMessagesTree(
					tree
						.addThreadPath(ctx.state.threadPath)
						.addList(ctx.state.tree, ctx.state.tree.threadID)
				)
				.addNewMessageForm(
					messageForm
						.addFilesList([])
						.addThreadName(ctx.state.parentName)
						.addThreadID(ctx.state.tree.threadID)
				)
				.addMessageHeader(
					messageHeader.addMessagePath(ctx.state.messagePath)
				)
				.addControlPanel(
					panel.addAuth(auth.addLogout())
				)
				.addHeader(
					header
						.addNewNote(ctx.state.newNoteHref)
						.addSearch(searchForm)
				)
		)

	ctx.body = layout.build()
	ctx.status = 200

	ctx.log.info("<-- newThreadMessage")
}

export default newThreadMessage
