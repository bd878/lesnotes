import MainBuilder from '../builders/mainBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import SettingsBuilder from '../builders/settingsBuilder';
import SidebarBuilder from '../builders/sidebarBuilder';

async function main(ctx) {
	console.log("--> main")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MainBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const sidebar = new SidebarBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	sidebar.addSettings(settings)

	content.addAuthorization()
	content.addSidebar(sidebar)

	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- main")
}

export default main;
