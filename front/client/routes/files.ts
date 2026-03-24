import * as is from '../third_party/is';
import FilesBuilder from '../builders/filesBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function files(ctx) {
	console.log("--> files")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new FilesBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	content.addNavigation()
	content.addControlPanel()
	content.addMessagesTree(ctx.state.stack)
	content.addFilesForm()
	if (is.notEmpty(ctx.state.files)) {
		content.addFilesList(ctx.state.files.files)
		content.addPagination(ctx.state.files.paging)
	}
	content.addLogout()

	header.addSearch()

	layout.addSettings(settings)
	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- files")
}

export default files
