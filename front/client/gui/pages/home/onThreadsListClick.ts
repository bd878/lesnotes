import * as is from '../../../third_party/is';
import {openThread, showMessage} from './messageCommands'

function onThreadsListClick(_elems, e) {
	if (is.notUndef(e.target.dataset.threadId)) {
		openThread(e.target.dataset.threadId)
	} else if (is.notEmpty(e.target.dataset.messageId)) {
		showMessage(e.target.dataset.messageId)
	}
}

export default onThreadsListClick
