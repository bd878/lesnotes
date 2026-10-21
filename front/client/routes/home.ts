import * as is from '../third_party/is';
import newMessage from './newMessage';
import homeTree from './homeTree';

async function home(ctx) {
	ctx.log.info("--> home")

	if (ctx.userAgent.isMobile) {
		if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0) {
			await homeTree(ctx)
		} else {
			await newMessage(ctx)
		}
	} else {
		await newMessage(ctx)
	}

	ctx.log.info("<-- home")
}

export default home
