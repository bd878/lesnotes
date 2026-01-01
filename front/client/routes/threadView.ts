import ThreadViewBuilder from '../builders/threadViewBuilder';

async function threadView(ctx) {
	console.log("--> threadView")

	const builder = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addThreadView(ctx.state.thread)
	await builder.addSettings()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- threadView")
}

export default threadView;
