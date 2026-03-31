import * as is from '../third_party/is';
import PublicMessageBuilder from '../builders/publicMessageBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import AuthBuilder from '../builders/authBuilder';

async function publicMessage(ctx) {
	console.log("--> publicMessage")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new PublicMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	ctx.state.messageFeatures.addNavigation(ctx.state.messageNavigation)

	content.addMessageView(ctx.state.me.ID, ctx.state.message)

	if (ctx.state.isAuthed) {
		auth.addLogout()
	} else {
		auth.addSignup()
	}

	header.addAuth(auth)
	content.addHeader(header)
	content.addMessageFeatures(ctx.state.messageFeatures)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- publicMessage")
}

export default publicMessage;
