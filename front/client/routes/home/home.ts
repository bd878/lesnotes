import HomeBuilder from './builder'

async function home(ctx) {
	console.log("--> home")

	const builder = new HomeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.search, ctx.path)

	await builder.addNewMessageForm(ctx.state.thread)
	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addMessagesList(ctx.state.stack)
	await builder.addFilesList()
	await builder.addFilesForm()
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(undefined, ctx.state.theme, ctx.state.fontSize, false)
	ctx.status = 200;

	console.log("<-- home")
}

export default home;
