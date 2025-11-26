import type {SearchMessage} from './models';
import api from './api';
import readPathJson from './readPathJson';
import models from './models';

async function searchMessagesPathJson(token: string, messages: SearchMessage[]) {
	let result = {
		error:    models.error(),
		messages: [],
	}

	try {
		for (const message of messages) {
			const path = await readPathJson(token, message.ID)

			path.path.reverse()
			result.messages.push(models.searchMessagePath(message, path.threadID, path.path))
		}
	} catch (e) {
		result.error.error    = true
		result.error.status   = 500
		result.error.explain  = e.toString()
	}

	return result
}

export default searchMessagesPathJson;
