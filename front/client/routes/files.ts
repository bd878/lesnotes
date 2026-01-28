import * as is from '../third_party/is';
import FilesBuilder from '../builders/filesBuilder';

async function files(ctx) {
	console.log("--> files")

	const builder = new FilesBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	await builder.addSettings()
	await builder.addNavigation()
	await builder.addControlPanel()
	await builder.addMessagesStack(ctx.state.stack)
	await builder.addFilesForm()
	if (is.notEmpty(ctx.state.files)) {
		await builder.addFilesList(ctx.state.files.files)
		await builder.addPagination(ctx.state.files.paging)
	}
	await builder.addSearch() // TODO: addFilesSearch
	await builder.addLogout()
	await builder.addSidebar()
	await builder.addFooter()

	ctx.body = await builder.build()
	ctx.status = 200;

	console.log("<-- files")
}

export default files
