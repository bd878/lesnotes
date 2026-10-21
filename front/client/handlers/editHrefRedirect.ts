import * as is from '../third_party/is'

async function editHrefRedirect(ctx, next) {
	ctx.log.info("--> editHrefRedirect")

	if (is.notEmpty(ctx.params.parentName)) {
		ctx.state.editHrefRedirect = "/" + ctx.params.parentName + "/" + ctx.params.messageName + ctx.search
	} else if (is.notEmpty(ctx.params.messageOrParentName)) {
		ctx.state.editHrefRedirect = "/" + ctx.params.messageOrParentName + ctx.search
	} else {
		ctx.state.editHrefRedirect = ""
	}

	await next()

	ctx.log.info("<-- editHrefRedirect")
}

export default editHrefRedirect
