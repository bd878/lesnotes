import LayoutBuilder from '../builders/layoutBuilder';
import AuthBuilder from '../builders/authBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SearchBuilder from '../builders/searchBuilder';
import ControlPanelBuilder from '../builders/controlPanelBuilder';

async function searchResult(ctx) {
	ctx.log.info("--> searchResult")

	const panel = new ControlPanelBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const auth = new AuthBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const content = new SearchBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)

	try {
		layout
			.addFooter()
			.addContent(
				content
					.addHeader(
						ctx.state.isAuthed
							? header.addSearch().addNewNote("/home" + ctx.search)
							: header.addSearch()
					)
					.addControlPanel(panel.addAuth(ctx.state.isAuthed ? auth.addLogout() : auth.addLogin()))
					.addMessagesList(ctx.state.list)
			)

		ctx.body = layout.build()
		ctx.status = 200;
	} catch (e) {
		console.log(e)
	}

	ctx.log.info("<-- searchResult")
}

export default searchResult
