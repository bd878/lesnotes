import PublicThreadBuilder from '../builders/publicThreadBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';

async function publicThread(ctx) {
	console.log("--> publicThread")

	const content = new PublicThreadBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	if (ctx.state.isAuthed) {
		auth.addLogout()
	} else {
		auth.addSignup()
	}

	header.addAuth(auth)
	content.addHeader(header)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
