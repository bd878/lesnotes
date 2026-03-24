import SignupBuilder from '../builders/signupBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function signup(ctx) {
	console.log("--> signup")

	const content = new SignupBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	content.addUsername()
	content.addPassword()
	content.addSubmit()

	header.addSearch()

	layout.addSettings(settings)
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- signup")
}

export default signup;
