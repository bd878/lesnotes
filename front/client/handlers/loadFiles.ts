import listFilesJson, { EmptyFilesList } from '../api/listFilesJson';
import * as is from '../third_party/is';

// TODO: load all files without limit, or imagine how to make a pagination
const limit = 10_000

async function loadFiles(ctx, next) {
	console.log("--> loadFiles")

	const userID = parseInt(ctx.params.user) || 0;
	const token = ctx.state.token
	const params = new URLSearchParams(ctx.search)
	const offset = parseInt(params.get("files")) || 0;

	if (is.notEmpty(token)) {
		// all user files
		ctx.state.files = await listFilesJson(token, 0, limit, offset)
	} else if (is.notEmpty(userID)) {
		// public user files
		ctx.state.files = await listFilesJson('', userID, limit, offset)
	}

	if (is.notEmpty(ctx.state.files)) {
		if (ctx.state.files.error.error) {
			console.error(ctx.state.files.error)
			ctx.body = "error"
			ctx.status = 400
			return
		}

		ctx.state.files.files.reverse()
	} else {
		ctx.state.files = EmptyFilesList
	}

	await next()

	console.log("<-- loadFiles")
}

export default loadFiles
