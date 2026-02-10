import PublicTranslationBuilder from '../builders/publicTranslationBuilder';

async function publicTranslation(ctx) {
	console.log("--> publicTranslation")

	const builder = new PublicTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addTranslationView()
	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200;

	console.log("<-- publicTranslation")
}

export default publicTranslation;
