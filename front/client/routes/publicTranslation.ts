import PublicTranslationBuilder from '../builders/publicTranslationBuilder';

async function publicTranslation(ctx) {
	console.log("--> publicTranslation")

	const builder = new PublicTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addTranslations(ctx.state.messageName, ctx.state.message.translations)
	await builder.addTranslationView(ctx.state.translation)
	await builder.addFilesView(ctx.state.message.files)
	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.translation)
	ctx.status = 200;

	console.log("<-- publicTranslation")
}

export default publicTranslation;
