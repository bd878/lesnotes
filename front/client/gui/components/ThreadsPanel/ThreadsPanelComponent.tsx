import React from 'react'
import MessagesList from '../MessagesList';
import MainMessage from '../MainMessage';

function ThreadsPanelComponent(props) {
	const {
		threadMessage,
		messages,
	} = props

	return (
		<div>
			<MainMessage message={threadMessage} />
			<MessagesList messages={messages} 	/>
		</div>
	)
}

export default ThreadsPanelComponent;