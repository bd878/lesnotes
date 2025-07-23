import React, {useEffect} from 'react';
import ReactDOM from 'react-dom/client';
import {connect} from '../../../third_party/react-redux';
import {
	LIMIT_LOAD_BY,
	LOAD_ORDER,
} from './const';
import {
	fetchMessagesActionCreator,
	selectMessages,
	selectThreadID,
	selectError,
} from '../../features/stack';
import Tag from '../Tag';
import List from '../List';
import MiniappListElement from '../MiniappListElement'

function MiniappThreadContainer(props) {
	const {
		threadID,
		messages,
		error,
		fetchMessages,
	} = props

	useEffect(() => {
		fetchMessages({limit: LIMIT_LOAD_BY, offset: 0, order: LOAD_ORDER, threadID: threadID})
	}, [fetchMessages, threadID]);

	return (
		<Tag>
			{error ? null : (
				<List el="ul" css="w-full flex flex-col space-y-px -mb-2">
					{messages.map((message) => (
						<MiniappListElement
							key={message.ID}
							message={message}
							margin="mb-2"
							radius="rounded-sm"
							textColor="text-main"
						/>
					))}
				</List>
			)}
		</Tag>
	)
}

const mapStateToProps = (state, {index}) => ({
	messages: selectMessages(index)(state),
	threadID: selectThreadID(index)(state),
	error: selectError(index)(state),
})

const mapDispatchToProps = (dispatch, {index}) => ({
	fetchMessages: payload => dispatch(fetchMessagesActionCreator(index)(payload)),
})

export default connect(mapStateToProps, mapDispatchToProps)(MiniappThreadContainer);
