import HomeBuilder from './builder'

async function home(ctx) {
	console.log("--> home")

	const builder = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.search)

	await builder.addNewMessageForm()
	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addMessagesList(ctx.state.stack)
	await builder.addFilesList(ctx.state.message)
	await builder.addFilesForm(ctx.state.message)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message, ctx.state.theme, ctx.state.fontSize)
	ctx.status = 200;

	console.log("<-- home")
}

export default home;
