import React, {useEffect} from 'react';
import ReactDOM from 'react-dom/client';
import {
	LIMIT_LOAD_BY,
	LOAD_ORDER,
} from './const';
import {
	fetchMessagesActionCreator,
	selectMessages,
	selectThreadID,
} from '../../features/stack';
import Tag from '../../components/Tag';

function MiniappThreadContainer(props) {
	const {
		threadID,
		messages,
		fetchMessages,
	} = props

	useEffect(() => {
		fetchMessages({limit: LIMIT_LOAD_BY, offset: 0, order: LOAD_ORDER, threadID: threadID})
	}, [fetchMessages, threadID]);

	return (
		<Tag>
			{messages.map((message) => (
				<Tag key={message.ID} css="px-2 py-1">{message.text}</Tag>
			))}
		</Tag>
	)
}

const mapStateToProps = (state, {index}) => ({
	messages: selectMessages(index)(state),
	threadID: selectThreadID(index)(state),
})

const mapDispatchToProps = (dispatch, {index}) => ({
	fetchMessages: payload => dispatch(fetchMessagesActionCreator(index)(payload)),
})

export default connect(mapStateToProps, mapDispatchToProps)(MiniappThreadContainer);
