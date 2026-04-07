import SignupBuilder from '../builders/signupBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import SidebarBuilder from '../builders/sidebarBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function signup(ctx) {
	console.log("--> signup")

	const content = new SignupBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const sidebar = new SidebarBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	layout
		.addFooter()
		.addContent(
			content
				.addUsername()
				.addPassword()
				.addSubmit()
				.addSidebar(sidebar.addSettings(settings))
		)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- signup")
}

export default signup;
