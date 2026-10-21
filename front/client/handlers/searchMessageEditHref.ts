import * as is from '../third_party/is'

async function searchMessageEditHref(ctx, next) {
	ctx.log.info("--> messageEditHref")

	if (is.notEmpty(ctx.state.message.name)) {
		ctx.state.messageEditHref = "/" + ctx.state.message.name + "/edit" + ctx.search
	}

	await next()

	ctx.log.info("<-- messageEditHref")
}

export default searchMessageEditHref;
