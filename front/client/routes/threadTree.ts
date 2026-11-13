import * as is from '../third_party/is';
import mobileTree from './mobileTree';
import publicMessage from './publicMessage'
import newThreadMessage from './newThreadMessage';

async function threadTree(ctx) {
	ctx.log.info("--> threadTree")

	if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0 && ctx.userAgent.isMobile) {
		await mobileTree(ctx)
	} else {
		if (ctx.state.isAuthed && ctx.userAgent.isMobile) {
			await newThreadMessage(ctx)
		} else {
			await publicMessage(ctx)
		}
	}

	ctx.log.info("<-- threadTree")
}

export default threadTree
