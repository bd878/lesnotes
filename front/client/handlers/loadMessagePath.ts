import type { Message } from '../api/models';
import * as is from '../third_party/is';
import crop from '../utils/crop';
import readPathJson from '../api/readPathJson';

async function loadMessagePath(ctx, next) {
	const token = ctx.state.token
	const id = ctx.state.messageID

	console.log("--> loadMessagePath")

	if (is.notEmpty(token) && is.notEmpty(id)) {
		const result = await readPathJson(token, id)

		if (is.notEmpty(result)) {
			if (result.error.error) {
				console.error(result.error)
				ctx.body = "error"
				ctx.status = 400;
				return
			}

			ctx.state.messagePath = composePath(result.path)
		} else {
			ctx.state.messagePath = ""
		}
	} else {
		ctx.state.messagePath = ""
	}

	await next()

	console.log("<-- loadMessagePath")
}

export default loadMessagePath

function composePath(path: Message[]): string {
	let result = "/"

	for (let i = path.length-1; i >= 0; i--) {
		let msg = path[i]
		if (is.notEmpty(msg.title)) {
			result += crop(msg.title, 15)
		} else if (is.notEmpty(msg.text)) {
			result += crop(msg.text, 15)
		} else {
			result += msg.ID
		}
		result += "/"
	}

	result = result.slice(0, -1)

	return result
}
