import api from '../api';

const limit = parseInt(LIMIT)

async function loadStack(ctx, next) {
	console.log("--> loadStack")

	const id = parseInt(ctx.query.id) || 0
	const threadID = parseInt(ctx.query.cwd) || 0

	const offsets = buildThreadOffsets(new URLSearchParams(ctx.request.search))

	ctx.state.stack = await api.readStackJson(ctx.state.token, threadID, limit, offsets)

	for (const thread of ctx.state.stack.stack) {
		thread.isCenter = function() { return this.ID == thread.centerID }
		thread.isSelected = function() { return this.ID == id }
	}

	await next()

	console.log("<-- loadStack")
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
