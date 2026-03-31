import * as is from '../third_party/is';
import PublicThreadMessageBuilder from '../builders/publicThreadMessageBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const content = new PublicThreadMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.isAuthed) {
		auth.addLogout()
	} else {
		auth.addSignup()
	}

	content.addAuth(auth)

	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

export default publicThreadMessage
