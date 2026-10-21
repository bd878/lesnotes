import * as is from '../third_party/is';
import publicMessage from './publicMessage'
import publicThread from './publicThread'
import messageView from './messageView';

async function publicThreadOrMessage(ctx) {
	ctx.log.info("--> publicThreadOrMessage")

	console.log("publicThreadOrMessage")

	// message with child messages is thread
	if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0) {
		await publicThread(ctx)
	} else {
		if (ctx.state.isAuthed) {
			await messageView(ctx)
		} else {
			await publicMessage(ctx)
		}
	}

	ctx.log.info("<-- publicThreadOrMessage")
}

export default publicThreadOrMessage
