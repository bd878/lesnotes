import {editMessage} from './messageCommands'

function onMessageEditClick(elems, e) {
	e.stopPropagation()
	editMessage(parseInt(elems.messageEditElem.dataset.messageId))
}

export default onMessageEditClick
