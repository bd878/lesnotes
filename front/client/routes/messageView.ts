import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	if (ctx.state.msg == "comments") {
		content.addMessageNavigation()
		content.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		content.addMessageNavigation()
		content.addFilesView(ctx.state.message.files)
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			content.addMessageNavigation()
			content.addFilesView(ctx.state.message.files)
		} else {
			content.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	content.addMessagesTree(ctx.state.stack)
	content.addMessageView(ctx.state.me.ID, ctx.state.message)
	content.addNewTranslation(ctx.state.message.ID)
	content.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	content.addLogout()

	header.addSearch()

	layout.addSettings()
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- messageView")
}

export default messageView;
