import LoginBuilder from '../builders/loginBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import SettingsBuilder from '../builders/settingsBuilder';
import SidebarBuilder from '../builders/sidebarBuilder';

async function login(ctx) {
	console.log("--> login")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new LoginBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const sidebar = new SidebarBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

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

	console.log("<-- login")
}

export default login;
