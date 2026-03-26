import * as is from '../third_party/is';
import MessageEditViewBuilder from '../builders/messageEditViewBuilder';
import MessagesTreeBuilder from '../builders/messagesTreeBuilder';
import LayoutBuilder from '../builders/layoutBuilder';
import HeaderBuilder from '../builders/headerBuilder';
import SettingsBuilder from '../builders/settingsBuilder';

async function messageEdit(ctx) {
	console.log("--> messageEdit")

	const layout = new LayoutBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const content = new MessageEditViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path)
	const header = new HeaderBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const settings = new SettingsBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);
	const tree = new MessagesTreeBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.state.theme, ctx.state.fontSize, ctx.search, ctx.path);

	header.addNewNote()
	tree.addList(ctx.state.tree)

	content.addMessagesTree(tree)
	content.addFilesSelector(ctx.state.files.files)
	content.addMessageEditForm(ctx.state.message)
	content.addNewTranslation(ctx.state.message.ID)
	content.addTranslations(ctx.state.message.ID, ctx.state.message.translations)

	layout.addFooter()
	layout.addHeader(header)
	layout.addContent(content)

	ctx.body = layout.build()
	ctx.status = 200

	console.log("<-- messageEdit")
}

export default messageEdit;
