import * as is from '../third_party/is';
import PublicThreadMessageBuilder from '../builders/publicThreadMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const content = new PublicThreadMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.msg == "comments") {
		content.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		content.addFilesView(ctx.state.message.files)
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			content.addFilesView(ctx.state.message.files)
		} else {
			content.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	content.addMessagesList(ctx.params.threadName /* TODO: use from load_path, it is message now, thread required */, ctx.state.messages)
	content.addMessageNavigation(messageNavigation)
	content.addTranslations(ctx.state.messageName, ctx.state.threadName, ctx.state.message.translations)
	content.addMessageView(ctx.state.message)

	if (ctx.state.me.ID) {
		// header.addLogout(logout)
	} else {
		content.addSignup()
	}

	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

export default publicThreadMessage
