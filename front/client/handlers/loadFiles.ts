import api from '../api';
import * as is from '../third_party/is';

const limit = parseInt(LIMIT)

async function loadFiles(ctx, next) {
	console.log("--> loadFiles")

	const userID = parseInt(ctx.params.user) || 0;
	const token = ctx.state.token
	const params = new URLSearchParams(ctx.search)
	const offset = parseInt(params.get("files")) || 0;

	if (is.notEmpty(token)) {
		// all user files
		ctx.state.files = await api.listFilesJson(token, 0, limit, offset)
	} else if (is.notEmpty(userID)) {
		// public user files
		ctx.state.files = await api.listFilesJson('', userID, limit, offset)
	}

	if (is.notEmpty(ctx.state.files)) {
		if (ctx.state.files.error.error) {
			console.error(ctx.state.files.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}
	} else {
		ctx.state.files = undefined
	}

	await next()

	console.log("<-- loadFiles")
}

export default loadFiles
