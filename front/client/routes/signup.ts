import SignupBuilder from '../builders/signupBuilder'
import LayoutBuilder from '../builders/layoutBuilder';
import SidebarBuilder from '../builders/sidebarBuilder';

async function signup(ctx) {
	ctx.log.info("--> signup")

	const content = new SignupBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const sidebar = new SidebarBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	layout
		.addFooter()
		.addContent(
			content
				.addUsername()
				.addPassword()
				.addSubmit()
				.addSidebar(sidebar)
		)

	ctx.body = layout.build()
	ctx.status = 200;

	ctx.log.info("<-- signup")
}

export default signup;
