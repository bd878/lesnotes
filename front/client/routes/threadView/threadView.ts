import type { Message } from '../../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../../api';
import * as is from '../../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from '../home/builder';

async function threadView(ctx) {
	console.log("--> threadView")

	const builder = new ThreadViewBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.search, ctx.path)

	await builder.addThreadView()
	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addMessagesList(ctx.state.stack)
	await builder.addSearch()
	await builder.addSidebar(ctx.search)
	await builder.addFooter()

	ctx.body = await builder.build(ctx.state.message, ctx.state.theme, ctx.state.fontSize, false)
	ctx.status = 200

	console.log("<-- threadView")
}

class ThreadViewBuilder extends HomeBuilder {
	threadView = undefined;
	async addThreadView() {
		return "not implemented"
	}
}

export default threadView;
