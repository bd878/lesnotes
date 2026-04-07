import * as is from '../third_party/is';
import PublicTranslationBuilder from '../builders/publicTranslationBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import MessageNavigationBuilder from '../builders/messageNavigationBuilder';

async function publicTranslation(ctx) {
	console.log("--> publicTranslation")

	const content = new PublicTranslationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const messageNavigation = new MessageNavigationBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	layout.addFooter().addHeader(header).addContent(content)

	ctx.body = layout.build()
	ctx.status = 200;

	console.log("<-- publicTranslation")
}

export default publicTranslation;
