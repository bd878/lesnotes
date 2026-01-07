import NewMessageBuilder from '../builders/newMessageBuilder'

async function newMessage(ctx) {
	console.log("--> newMessage")

	const builder = new NewMessageBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addFilesInput(ctx.state.files.files)
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addNewMessageForm(ctx.state.thread)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200;

	console.log("<-- newMessage")
}

export default newMessage;
