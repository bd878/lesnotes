import type { Error, ThreadMessages } from './models';
import api from './index';
import models from './models';
import * as is from '../third_party/is'

interface StackResponse {
	error:     Error;
	stack:     ThreadMessages[];
}

async function readStackJson(token: string, messageID: number/*, lastMessageID: number*/, limit: number, offsets: Record<number, number> = {}): Promise<StackResponse> {
	let result = {
		error:       models.error(),
		stack:       [],
	}

	const path = await api.readPathJson(token, messageID)
	if (path.error.error) {
		result.error = path.error
		return result;
	}

	const ids = path.path.map(message => message.ID)
	ids.reverse()

	// len(messages) == len(path.path)
	const messageIDs = [0 /* root */ , ...ids] /* threadID == messageID */

	path.path.push(JSON.parse(JSON.stringify(models.EmptyMessage /* root message */)))
	path.path.reverse()

	// messageID = 0 : messages = [0], path.path = [EmptyMessage]
	for (
		let i = 0, isRoot = true, parentID = 0, messageID = 0 /* first is root : = 0 */, message = JSON.parse(JSON.stringify(models.EmptyMessage)) /* first is message : EmptyMessage */;
		i < messageIDs.length;
		i++, isRoot = false, parentID = messageIDs[i-1], messageID = messageIDs[i], message = path.path[i]
	) {
		const offset = is.undef(offsets[messageID]) ? 0 : offsets[messageID];

		const messages = await api.readMessagesJson(token, 0, messageID, 1 /* order */, limit, offset)
		if (messages.error.error) {
			result.error = messages.error
			return result;
		}

		messages.messages.reverse()

		const thread = JSON.parse(JSON.stringify(models.thread()))
		thread.parentID = parentID
		thread.isRoot = isRoot

		result.stack.push(models.threadMessages(messages.messages, messages.paging, message, thread))
	}

	return result;
}

export default readStackJson;
