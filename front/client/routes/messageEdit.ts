import MessageEditViewBuilder from '../builders/messageEditViewBuilder';

async function messageEdit(ctx) {
	console.log("--> messageEdit")

	const builder = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addFilesInput()
	await builder.addFilesList()
	await builder.addMessageEditForm(ctx.state.message)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- messageEdit")
}

export default messageEdit;
