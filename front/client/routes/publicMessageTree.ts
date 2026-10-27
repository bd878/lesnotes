import * as is from '../third_party/is';
import mobileTree from './mobileTree';
import publicMessage from './publicMessage'
import newThreadMessage from './newThreadMessage';

async function publicMessageTree(ctx) {
	ctx.log.info("--> publicMessageTree")

	if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0 && ctx.userAgent.isMobile) {
		await mobileTree(ctx)
	} else {
		if (ctx.state.isAuthed && ctx.userAgent.isMobile) {
			await newThreadMessage(ctx)
		} else {
			await publicMessage(ctx)
		}
	}

	ctx.log.info("<-- publicMessageTree")
}

export default publicMessageTree
