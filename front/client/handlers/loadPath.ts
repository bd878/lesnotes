import type { Message } from '../api/models';
import * as is from '../third_party/is';
import readPathJson from '../api/readPathJson';

async function loadPath(ctx, next) {
	const token = ctx.state.token
	const id = ctx.state.messageID

	console.log("--> loadPath")

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

	console.log("<-- loadPath")
}

export default loadPath

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

function crop(str: string, size: number): string {
	if (str.length > size) {
		return `${str.slice(0, size)}...`
	} else if (str.length <= size) {
		return str
	}
}