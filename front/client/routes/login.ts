import LoginBuilder from '../builders/loginBuilder';
import LayoutBuilder from '../builders/layoutBuilder';

async function login(ctx) {
	console.log("--> login")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new LoginBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	content.addUsername()
	content.addPassword()
	content.addSubmit()

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- login")
}

export default login;
