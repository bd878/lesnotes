import * as is from '../third_party/is';
import newRootMessage from './newRootMessage';
import homeTree from './homeTree';

async function home(ctx) {
	ctx.log.info("--> home")

	if (ctx.userAgent.isMobile) {
		if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0) {
			await homeTree(ctx)
		} else {
			await newRootMessage(ctx)
		}
	} else {
		await newRootMessage(ctx)
	}

	ctx.log.info("<-- home")
}

export default home
