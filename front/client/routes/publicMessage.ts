import * as is from '../third_party/is';
import PublicMessageBuilder from '../builders/publicMessageBuilder';
import MessageViewBuilder from '../builders/messageViewBuilder';
import PublicMessageHeaderBuilder from '../builders/publicMessageHeaderBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import AuthBuilder from '../builders/authBuilder';
import SearchFormBuilder from '../builders/searchFormBuilder';

async function publicMessage(ctx) {
	ctx.log.info("--> publicMessage")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const view = new MessageViewBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const messageHeader = new PublicMessageHeaderBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const searchForm = new SearchFormBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const content = new PublicMessageBuilder(ctx.state.isAuthed, ctx.state.messageName, ctx.state.messageName, ctx.userAgent.isMobile,
		ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	if (ctx.state.messageFeatures.hasFeatures()) {
		content.addMessageFeatures(ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation))
	}

	layout
		.addTitle(ctx.state.message.title || ctx.state.message.text)
		.addFooter()
		.addContent(
			content
				.addMessageHeader(
					messageHeader.addThreadLink(ctx.state.message.parentMessage)
				)
				.addMessageView(
					view
						.addDeleteRedirectUrl(ctx.state.deleteRedirectUrl)
						.addMessage(ctx.state.message)
				)
				.addHeader(ctx.state.isAuthed
					? header.addAuth(auth.addLogout()).addSearch(searchForm)
					: ctx.state.message.parentMessage.private
						? header.addAuth(auth.addLogin())
						: header.addAuth(auth.addLogin()).addSearch(searchForm.addThreadName(ctx.state.message.parentMessage.name))
				)
		)

	ctx.body = layout.build()
	ctx.status = 200;

	ctx.log.info("<-- publicMessage")
}

export default publicMessage;
