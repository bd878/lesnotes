import ThreadEditBuilder from '../builders/threadEditBuilder'

async function threadEdit(ctx) {
	console.log("--> threadEdit")

	const builder = new ThreadEditBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addThreadEditForm(ctx.state.thread)
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addSettings()
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(undefined, ctx.state.thread)
	ctx.status = 200

	console.log("<-- threadEdit")
}

export default threadEdit
