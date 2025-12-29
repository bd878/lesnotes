import type { Error, Message } from './models'
import api from './api';
import models from './models';

interface ReadPathResponse {
	error:    Error;
	path:     Message[];
	threadID: number;
}

// id = 0 : path = []
async function readPathJson(token: string, id: number): Promise<ReadPathResponse> {
	let result: ReadPathResponse = {
		error:       models.error(),
		path:        [],
		threadID:    0,
	}

	try {
		const [response, error] = await api('/messages/v2/read_path', {
			method: "POST",
			body: {
				token: token,
				req:   {
					id: id,
				},
			},
		});

		if (error.error) {
			result.error = models.error(error)
		} else {
			// do not .reverse() here; client may append NullThread and then reverse
			result.path     = response.path.map(models.message)
			result.threadID = response.thread
		}

	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result;
}

export default readPathJson;
