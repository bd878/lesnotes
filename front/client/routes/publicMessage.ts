import PublicMessageBuilder from '../builders/publicMessageBuilder';

async function publicMessage(ctx) {
	console.log("--> publicMessage")

	const builder = new PublicMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message)
	ctx.status = 200;

	console.log("<-- publicMessage")
}

export default publicMessage;
