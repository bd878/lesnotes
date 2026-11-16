import * as is from '../third_party/is'
import LoginBuilder from '../builders/loginBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import SidebarBuilder from '../builders/sidebarBuilder';

async function login(ctx) {
	ctx.log.info("--> login")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new LoginBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const sidebar = new SidebarBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	if (is.notEmpty(ctx.state.error)) {
		content.addError(ctx.state.error)
	}

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

	ctx.log.info("<-- login")
}

export default login;
