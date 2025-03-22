import React, {lazy} from 'react'
import Tag from '../Tag';

const MessageForm = lazy(() => import("../MessageForm"));
const MainMessage = lazy(() => import("../MainMessage"));

function MessagePageComponent(props) {
	const {
		message,
		sendMessage,
		updateMessage,
		resetEditMessage,
		messageForEdit,
	} = props

	return (
		<Tag>
			<MainMessage
				message={message}
			/>
			<MessageForm
				send={sendMessage}
				update={updateMessage}
				reset={resetEditMessage}
				messageForEdit={messageForEdit}
			/>
		</Tag>
	)
}

export default MessagePageComponent