import MainBuilder from '../builders/mainBuilder';

async function main(ctx) {
	console.log("--> main")

	const builder = new MainBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addFooter()
	await builder.addSidebar()
	await builder.addAuthorization()

	ctx.body = await builder.build()
	ctx.status = 200;

	console.log("<-- main")
}

export default main;
