import Config from 'config';
import HomeBuilder from './builder'

async function home(ctx) {
	console.log("--> home")

	const { me, stack, message, lang, theme, fontSize } = ctx.state

	ctx.set({ "Cache-Control": "no-cache,max-age=0" })

	const builder = new HomeBuilder(ctx.userAgent.isMobile, lang, ctx.search)

	await builder.addNewMessageForm()
	await builder.addSettings(undefined, lang, theme, fontSize)
	await builder.addMessagesList(undefined, stack)
	await builder.addFilesList(message)
	await builder.addFilesForm(message)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(message, theme, fontSize)
	ctx.status = 200;

	console.log("<-- home")

	return;
}

export default home;
