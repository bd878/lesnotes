import * as is from '../third_party/is';
import publicMessage from './publicMessage'
import publicThread from './publicThread'
import newThreadMessage from './newThreadMessage'
import messageView from './messageView';

async function publicThreadOrMessage(ctx) {
	ctx.log.info("--> publicThreadOrMessage")

	// message with child messages is thread
	if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0 && !ctx.userAgent.isMobile) {
		await publicThread(ctx)
	} else {
		if (ctx.state.isAuthed) {
			if (ctx.userAgent.isMobile) {
				await newThreadMessage(ctx)
			} else {
				await messageView(ctx)
			}
		} else {
			await publicMessage(ctx)
		}
	}

	ctx.log.info("<-- publicThreadOrMessage")
}

export default publicThreadOrMessage
