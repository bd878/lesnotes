import React, {useEffect, useCallback} from 'react';
import ThreadsPanelComponent from './ThreadsPanelComponent'
import {
	LIMIT_LOAD_BY,
	LOAD_ORDER,
} from './const';
import {
	fetchMessagesActionCreator,
	sendMessageActionCreator,
	resetActionCreator,
	selectMessages,
	selectThreadMessage,
	selectThreadID,
} from '../../features/threads';
import {connect} from '../../third_party/react-redux';
import * as is from '../../third_party/is';

function ThreadsPanelContainer(props) {
	const {
		fetch,
		threadMessage,
		reset,
		send,
		messages,
		threadID,
		shouldShowThreadsPanel,
	} = props

	useEffect(() => {
		if (threadID !== 0)
			fetch(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [threadID])

	const onResetClick = useCallback(reset, [reset])

	if (threadID == 0 && !shouldShowThreadsPanel)
		return null

	return (
		<ThreadsPanelComponent
			close={onResetClick}
			threadMessage={threadMessage}
			messages={messages}
			send={send}
		/>
	)
}

const mapStateToProps = state => ({
	threadMessage: selectThreadMessage(state),
	messages: selectMessages(state),
	threadID: selectThreadID(state),
})

const mapDispatchToProps = ({
	reset: resetActionCreator,
	fetch: fetchMessagesActionCreator,
	send: sendMessageActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadsPanelContainer);
