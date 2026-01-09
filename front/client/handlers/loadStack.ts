import type { ThreadMessages } from '../api/models'
import api from '../api';
import * as is from '../third_party/is';

interface OpenClosedThreadMessages extends ThreadMessages {
	isClose: boolean;
	openIDs: string;
	closeIDs: string;
}

const limit = parseInt(LIMIT)

async function loadStack(ctx, next) {
	console.log("--> loadStack")

	const id = parseInt(ctx.query.id) || 0
	const threadID = parseInt(ctx.query.cwd) || 0

	const params = new URLSearchParams(ctx.request.search)
	const offsets = buildThreadOffsets(params)

	ctx.state.stack = await api.readStackJson(ctx.state.token, threadID, limit, offsets)

	if (is.notEmpty(ctx.state.stack)) {
		if (ctx.state.stack.error.error) {
			console.error(ctx.state.stack.error)
			ctx.body = "error"
			ctx.status = 400;
			return;
		}
		ctx.state.stack = ctx.state.stack.stack.map(openClosed(closeIDs(params.get("close") || "")))
	} else {
		ctx.state.stack = []
	}

	for (const stack of ctx.state.stack) {
		stack.thread.isCenter = function() { return this.ID == stack.thread.centerID }
		stack.thread.isSelected = function() { return this.ID == id }
	}

	await next()

	console.log("<-- loadStack")
}

function closeIDs(closeStr: string = ""): number[] {
	return closeStr.split(",").map(parseFloat).filter(v => !isNaN(v))
}

type StackMapFn = (stack: ThreadMessages, index: number, arr: ThreadMessages[]) => OpenClosedThreadMessages;

function openClosed(closed: number[]): StackMapFn {
	return function(stack: ThreadMessages, index: number, arr: ThreadMessages[]): OpenClosedThreadMessages {
		let isClose = false;
		let openIDs = "";
		let closeIDs = "";

		const set = new Set(closed)

		if (set.has(stack.message.ID)) {
			isClose = true
			closeIDs = Array.from(set).join(",")
			set.delete(stack.message.ID)
			openIDs = Array.from(set).join(",")
		} else {
			isClose = false
			openIDs = Array.from(set).join(",")
			set.add(stack.message.ID)
			closeIDs = Array.from(set).join(",")
		}

		return {
			...stack,
			isClose,
			openIDs,
			closeIDs,
		}
	}
}

function buildThreadOffsets(searchParams): Record<number, number> {
	const threadToOffset = {}

	for (const [key, value] of searchParams) {
		const threadID = parseInt(key)
		const offset = parseInt(value)
		if (!isNaN(threadID) && !isNaN(offset))
			threadToOffset[threadID] = offset
	}

	return threadToOffset
}

export default loadStack
export type { OpenClosedThreadMessages };
