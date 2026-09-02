import Config from "config";
import { logs } from '@opentelemetry/api-logs';
import * as is from '../third_party/is'

const log = logs.getLogger('default');

let logger = data => {
	if (is.object(data)) {
		console.log(data.body)
	} else {
		console.log(data)
	}
}
if (Config.get("env") != "development") {
	logger = data => log.emit(data)
}

export default logger