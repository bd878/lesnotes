import React from 'react'
import Tag from '../Tag';
import MessagesList from '../MessagesList';
import MainMessage from '../MainMessage';
import MessageForm from '../MessageForm';

function ThreadsPanelComponent(props) {
	const {
		threadMessage,
		messages,
	} = props

	return (
		<Tag>
			<MainMessage message={threadMessage} />
			<MessagesList messages={messages} 	/>
			<MessageForm messageForEdit={{}} reset={() => {}} send={() => {}} update={() => {}} edit={() => {}} />
		</Tag>
	)
}

export default ThreadsPanelComponent;