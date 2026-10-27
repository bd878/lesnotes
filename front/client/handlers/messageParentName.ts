import * as is from '../third_party/is';

async function messageParentName(ctx, next) {
	ctx.log.info("--> messageParentName")

	if (is.notEmpty(ctx.state.message) && is.notEmpty(ctx.state.message.parentMessage.name)) {
		ctx.state.parentName = ctx.state.message.parentMessage.name
	}

	await next()

	ctx.log.info("<-- messageParentName")
}

export default messageParentName
