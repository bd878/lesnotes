import * as is from '../third_party/is';
import publicMessage from './publicMessage'
import publicThread from './publicThread'

async function publicThreadOrMessage(ctx) {
	console.log("--> publicThreadOrMessage")

	// message with child messages is thread
	if (is.notEmpty(ctx.state.tree) && ctx.state.tree.count > 0) {
		await publicThread(ctx)
	} else {
		// otherwise, just a message
		await publicMessage(ctx)
	}

	console.log("<-- publicThreadOrMessage")
}

export default publicThreadOrMessage
