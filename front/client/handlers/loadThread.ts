import api from '../api';
import * as is from '../third_party/is';

async function loadThread(ctx, next) {
	console.log("--> loadThread")

	await next()

	console.log("<-- loadThread")
}

export default loadThread
