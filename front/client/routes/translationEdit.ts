import TranslationEditViewBuilder from '../builders/translationEditViewBuilder';

async function translationEdit(ctx) {
	console.log("--> translationEdit")

	const builder = new TranslationEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addFilesView(ctx.state.message.files)
	await builder.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	await builder.addTranslationEditForm(ctx.state.message.ID, ctx.state.translation)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- translationEdit")
}

export default translationEdit;
