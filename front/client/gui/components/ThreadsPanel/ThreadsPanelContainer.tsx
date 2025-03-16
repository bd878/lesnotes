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
	selectThreadID,
} from '../../features/threads';
import {connect} from '../../third_party/react-redux';
import * as is from '../../third_party/is';

function ThreadsPanelContainer(props) {
	const {
		fetchMessages,
		threadMessage,
		messages,
		threadID,
		shouldShowThreadsPanel,
	} = props

	useEffect(() => {
		if (threadID != 0)
			fetchMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [threadID])

	if (threadID == 0 && !shouldShowThreadsPanel) {
		return null
	}

	return (
		<ThreadsPanelComponent
			threadMessage={threadMessage}
			messages={messages}
		/>
	)
}

const mapStateToProps = state => ({
	threadMessage: selectThreadMessage(state),
	messages: selectMessages(state),
	threadID: selectThreadID(state),
})

const mapDispatchToProps = ({
	fetchMessages: fetchMessagesActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadsPanelContainer);
