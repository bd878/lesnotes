import * as is from '../third_party/is';
import MessageEditViewBuilder from '../builders/messageEditViewBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';
import HomeBuilder from '../builders/homeBuilder';
import PublicMessagesTreeBuilder from '../builders/publicMessagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageHeaderBuilder from '../builders/messageHeaderBuilder';

async function editThreadMessage(ctx) {
	ctx.log.info("--> editThreadMessage")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const messageForm = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageHeader = new MessageHeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new PublicMessagesTreeBuilder(ctx.state.isAuthed, ctx.state.parentName, ctx.state.messageName,
		ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout
		.addTitle(ctx.state.message.title || ctx.state.message.text)
		.addFooter()
		.addContent(
			content
				.addMessagesTree(
					tree
						.addThreadPath(ctx.state.cwdPath)
						.addList(ctx.state.tree, ctx.state.tree.threadID)
				)
				.addMessageEditForm(
					messageForm
						.addRedirectUrl(ctx.path + ctx.search)
						.addCancelHref(ctx.state.cancelEditHref)
						.addFilesList(ctx.state.message.files)
						.addMessage(ctx.state.message)
				)
				.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
				.addMessageHeader(messageHeader.addMessagePath(ctx.state.messagePath))
				.addHeader(
					header
						.addNewNote(ctx.state.newNoteHref)
						.addSearch()
				)
				.addControlPanel(panel.addAuth(auth.addLogout()))
		)

	ctx.body = layout.build()
	ctx.status = 200

	ctx.log.info("<-- editThreadMessage")
}

export default editThreadMessage
