import * as is from '../third_party/is'

async function deleteRedirectUrl(ctx, next) {
	ctx.log.info("--> deleteRedirectUrl")

	if (is.notEmpty(ctx.params.parentName)) {
		ctx.state.deleteRedirectUrl = "/" + ctx.params.parentName + "/new" + ctx.search
	} else if (is.notEmpty(ctx.params.messageOrParentName)) {
		ctx.state.deleteRedirectUrl = "/" + ctx.state.message.parentMessage.name + "/new" + ctx.search
	} else {
		ctx.state.deleteRedirectUrl = ""
	}

	await next()

	ctx.log.info("<-- deleteRedirectUrl")
}

export default deleteRedirectUrl
