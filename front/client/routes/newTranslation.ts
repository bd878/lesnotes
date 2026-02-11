import NewTranslationBuilder from '../builders/newTranslationBuilder'

async function newTranslation(ctx) {
	console.log("--> newTranslation")

	const builder = new NewTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addNewTranslationForm(ctx.state.messageID)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200;

	console.log("<-- newTranslation")
}

export default newTranslation;
