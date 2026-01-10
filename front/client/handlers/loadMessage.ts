import type { Message, File } from '../api/models'
import type { FileWithMime } from '../types'
import readMessageJson, { EmptyReadMessage } from '../api/readMessageJson';
import * as is from '../third_party/is';

interface MessageWithFilesMime extends Message {
	files: FileWithMime[];
}

async function loadMessage(ctx, next) {
	const id = parseInt(ctx.query.id) || parseInt(ctx.params.id) || 0
	const name = ctx.params.messageName || ""
	const userID = parseInt(ctx.params.user) || 0
	const token = ctx.state.token

	console.log("--> loadMessage")

	if (is.notEmpty(token)) {
		if (is.notEmpty(id)) {
			ctx.state.message = await readMessageJson(token, 0 /* me */, id)
		} else if (is.notEmpty(name)) {
			ctx.state.message = await readMessageJson(token, 0 /* me */, 0, name /* public name */)
		}
	} else if (is.notEmpty(userID)) {
		if (is.notEmpty(id)) {
			ctx.state.message = await readMessageJson("", userID, id)
		}
	} else if (is.notEmpty(name)) {
		ctx.state.message = await readMessageJson("", 0 /* public */, 0, name /* public name */)
	}

	if (is.notEmpty(ctx.state.message)) {
		if (ctx.state.message.error.error) {
			console.error(ctx.state.message.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}

		ctx.state.message = ctx.state.message.message
		ctx.state.message.files = ctx.state.message.files.map(fileWithMime)
	} else {
		ctx.state.message = EmptyReadMessage.message
	}

	await next()

	console.log("<-- loadMessage")
}

export default loadMessage
export type { FileWithMime }

function fileWithMime(file: File, index: number, arr: File[]): FileWithMime {
	const result = {
		...file,
		isDocument: false,
		isImage:    false,
		isAudio:    false,
		isVideo:    false,
		isText:     false,
		isFile:     false,
	}

	if (file.mime.includes("image")) {
		result.isImage = true
	} else if (file.mime.includes("pdf")) {
		result.isDocument = true
	} else {
		result.isFile = true
	}

	return result
}