import * as is from '../third_party/is'

async function messageEditHref(ctx, next) {
	ctx.log.info("--> messageEditHref")

	if (is.notEmpty(ctx.state.parentName)) {
		ctx.state.messageEditHref = ctx.path + "/edit" + ctx.search
	}

	await next()

	ctx.log.info("<-- messageEditHref")
}

export default messageEditHref
