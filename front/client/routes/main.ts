import MainBuilder from '../builders/mainBuilder';
import LayoutBuilder from '../builders/layoutBuilder';

async function main(ctx) {
	console.log("--> main")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MainBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	content.addAuthorization()

	layout.addSettings()
	layout.addFooter()
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- main")
}

export default main;
