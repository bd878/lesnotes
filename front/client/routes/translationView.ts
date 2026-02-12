import TranslationViewBuilder from '../builders/translationViewBuilder';

async function translationView(ctx) {
	console.log("--> translationView")

	const builder = new TranslationViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addFilesView(ctx.state.message.files)
	await builder.addNewTranslation(ctx.state.message.ID)
	await builder.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	await builder.addTranslationView(ctx.state.message.ID, ctx.state.translation)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- translationView")
}

export default translationView;
