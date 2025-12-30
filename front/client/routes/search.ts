import SearchBuilder from '../builders/searchBuilder'

async function search(ctx) {
	console.log("--> search")

	const builder = new SearchBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addMessagesList(ctx.state.messages)
	await builder.addFilesList()
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- search")
}

export default search
