import models from './models';
import api from './api';

async function sendLog(body: any) {
	let result = {
		error:       models.error(),
	}

	try {
		let [response, error] = await api("/telemetry/v1/send", {
			method:    "POST",
			isFullUrl: true,
			body:      body,
		});

		if (error)
			result.error = models.error(error)
	} catch (e) {
		result.error.error   = true
		result.error.status  = 500
		result.error.explain = e.toString()
	}

	return result
}

export default sendLog