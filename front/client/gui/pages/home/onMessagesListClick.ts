import * as is from '../../../third_party/is';
import {showMessage, paginateMessages, openThread} from './messageCommands'

function onMessagesListClick(elems, e) {
	if (is.notEmpty(e.target.dataset.messageId)) {
		showMessage(e.target.dataset.messageId)
	} else if (is.notEmpty(e.target.dataset.threadId) && is.notEmpty(e.target.dataset.direction)) {
		paginateMessages(e.target.dataset.threadId, e.target.dataset.direction, e.target.dataset.offset, e.target.dataset.limit)
	} else if (is.notEmpty(e.target.dataset.threadId)) {
		openThread(e.target.dataset.threadId)
	}
}

export default onMessagesListClick
