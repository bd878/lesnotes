import PublicThreadMessageBuilder from '../builders/publicThreadMessageBuilder'

async function publicThreadMessage(ctx) {
	console.log("--> publicThreadMessage")

	const builder = new PublicThreadMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessagesList(ctx.params.threadName /* TODO: use from load_path, it is message now, thread required */, ctx.state.messages)
	await builder.addSearch()
	await builder.addTranslations(ctx.state.messageName, ctx.state.threadName, ctx.state.message.translations)
	await builder.addMessageView(ctx.state.message)
	await builder.addFilesView(ctx.state.message.files)
	await builder.addSettings()

	if (ctx.state.me.ID) {
		await builder.addLogout()
	} else {
		await builder.addSignup()
	}

	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message)
	ctx.status = 200

	console.log("<-- publicThreadMessage")
}

export default publicThreadMessage
