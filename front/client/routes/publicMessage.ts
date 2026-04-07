import * as is from '../third_party/is';
import PublicMessageBuilder from '../builders/publicMessageBuilder';
import MessageViewBuilder from '../builders/messageViewBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import AuthBuilder from '../builders/authBuilder';

async function publicMessage(ctx) {
	console.log("--> publicMessage")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const content = new PublicMessageBuilder(ctx.state.isAuthed, ctx.state.threadName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	layout
		.addFooter()
		.addContent(
			content
				.addMessageView(
					view
						.addDeleteRedirectUrl("/t/" + ctx.state.threadName + ctx.search)
						.addMessage(ctx.state.message)
				)
				.addHeader(
					header.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin())
				)
				.addMessageFeatures(
					ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)
				)
		)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- publicMessage")
}

export default publicMessage;
