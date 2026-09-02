import api from '../api';
import readMessageJson, { EmptyReadMessage } from '../api/readMessageJson';
import * as is from '../third_party/is';

async function loadParentMessage(ctx, next) {
	const token = ctx.state.token
	const userID = parseInt(ctx.params.user) || 0
	const id = ctx.state.threadID
	const name = ctx.state.parentName || ""

	ctx.log.info("--> loadParentMessage")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.parentMessage = await readMessageJson(token, 0 /* me */, id)
		} else if (is.notEmpty(name)) {
			ctx.state.parentMessage = await readMessageJson(token, 0 /* me */, 0, name /* public name */)
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			ctx.state.parentMessage = await readMessageJson("", userID, id)
		}
	} else if (is.notEmpty(name)) {
		ctx.state.parentMessage = await readMessageJson("", 0 /* public */, 0, name /* public name */)
	}

	if (is.notEmpty(ctx.state.parentMessage)) {
		if (ctx.state.parentMessage.error.error) {
			ctx.log.error(ctx.state.parentMessage.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}

		ctx.state.parentMessage = ctx.state.parentMessage.message
		// TODO: see loadMessage : fileWithMime, sortImagesFirst
	} else {
		ctx.state.parentMessage = EmptyReadMessage.message
	}

	await next()

	ctx.log.info("<-- loadParentMessage")
}

export default loadParentMessage
