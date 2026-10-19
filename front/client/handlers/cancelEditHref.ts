import * as is from '../third_party/is'

async function cancelEditHref(ctx, next) {
	ctx.log.info("--> cancelEditHref")

	if (is.notEmpty(ctx.params.parentName)) {
		ctx.state.cancelEditHref = "/" + ctx.params.parentName + "/" + ctx.params.messageName + ctx.search
	} else if (is.notEmpty(ctx.params.messageOrParentName)) {
		ctx.state.cancelEditHref = "/" + ctx.params.messageOrParentName + ctx.search
	} else {
		ctx.state.cancelEditHref = ""
	}

	console.log(ctx.state.cancelEditHref)

	await next()

	ctx.log.info("<-- cancelEditHref")
}

export default cancelEditHref
