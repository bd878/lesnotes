import * as is from '../third_party/is';
import MessageViewBuilder from '../builders/messageViewBuilder';

async function messageView(ctx) {
	console.log("--> messageView")

	const builder = new MessageViewBuilder(ctx.userAgent.isMobile, ctx.state.lang,
		ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()

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

	await builder.addMessagesStack(ctx.state.stack)
	await builder.addMessageView(ctx.state.me.ID, ctx.state.message)
	await builder.addNewTranslation(ctx.state.message.ID)
	await builder.addTranslations(ctx.state.message.ID, ctx.state.message.translations)
	await builder.addSearch()
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200

	console.log("<-- messageView")
}

export default messageView;
