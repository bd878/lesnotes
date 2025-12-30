import PublicThreadBuilder from '../builders/publicThreadBuilder'

async function publicThread(ctx) {
	console.log("--> publicThread")

	const builder = new PublicThreadBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessagesList(ctx.params.threadName /* TODO: use from load_path, ctx.thread.name is message name now, but thread name required */, ctx.state.messages)
	await builder.addSearch()
	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- publicThread")
}

export default publicThread
