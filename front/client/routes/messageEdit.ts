import MessageEditViewBuilder from '../builders/messageEditViewBuilder';

async function messageEdit(ctx) {
	console.log("--> messageEdit")

	const builder = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()
	await builder.addFilesView(ctx.state.message.files)
	await builder.addFilesSelector(ctx.state.files.files)
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addMessageEditForm(ctx.state.message)
	await builder.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- messageEdit")
}

export default messageEdit;
