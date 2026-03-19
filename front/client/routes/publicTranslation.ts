import * as is from '../third_party/is';
import PublicTranslationBuilder from '../builders/publicTranslationBuilder';

async function publicTranslation(ctx) {
	console.log("--> publicTranslation")

	const builder = new PublicTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	if (ctx.state.msg == "comments") {
		await builder.addMessageNavigation()
		await builder.addComments(ctx.state.message.ID, ctx.state.comments)
	} else if (ctx.state.msg == "files") {
		await builder.addMessageNavigation()
		await builder.addFilesView(ctx.state.message.files)
	} else {
		if (is.array(ctx.state.message.files) && ctx.state.message.files.length > 0) {
			await builder.addMessageNavigation()
			await builder.addFilesView(ctx.state.message.files)
		} else {
			await builder.addComments(ctx.state.message.ID, ctx.state.comments)
		}
	}

	await builder.addTranslations(ctx.state.messageName, ctx.state.message.translations)
	await builder.addTranslationView(ctx.state.translation)
	await builder.addSettings()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.translation)
	ctx.status = 200;

	console.log("<-- publicTranslation")
}

export default publicTranslation;
