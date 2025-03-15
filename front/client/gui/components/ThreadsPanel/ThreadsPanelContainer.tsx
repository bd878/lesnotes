import React, {useEffect} from 'react';
import ThreadsPanelComponent from './ThreadsPanelComponent'
import {
  LIMIT_LOAD_BY,
  LOAD_ORDER,
} from './const';
import {
	fetchMessagesActionCreator,
	selectMessages,
	selectThreadMessage,
} from '../../features/threads';
import {connect} from '../../third_party/react-redux';

function ThreadsPanelContainer(props) {
	const {
		fetchMessages,
		threadMessage,
		messages,
	} = props

	useEffect(() => {
		fetchMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [])

	return (
		<ThreadsPanelComponent
			threadMessage={threadMessage}
			messages={messages}
		/>
	)
}

const mapStateToProps = state => ({
	threadMessage: selectThreadMessage(state),
	messages: selectMessages(state)
})

const mapDispatchToProps = ({
	fetchMessages: fetchMessagesActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadsPanelContainer);
