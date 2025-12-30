import MessageViewBuilder from '../builders/messageViewBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const builder = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addMessageView(ctx.state.me.ID, ctx.state.message)
	await builder.addSettings()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addSearch()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message)
	ctx.status = 200

	console.log("<-- messageView")
}

export default messageView;
