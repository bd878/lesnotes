import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {
	setEditMessageActionCreator,
	deleteMessageActionCreator,
	copyMessageActionCreator,
} from '../../features/messages';
import {
	setThreadMessageActionCreator,
	selectThreadID,
} from '../../features/threads';
import {connect} from '../../../third_party/react-redux';

function MessagesListContainer(props) {
	const {
		css,
		messages,
		loading,
		error,
		threadID,
		setThreadMessage,
		setEditMessage,
		deleteMessage,
		copyMessage,
	} = props

	const onToggleThreadClick = useCallback(setThreadMessage, [setThreadMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback(copyMessage, [copyMessage])

	return (
		<MessagesListComponent
			css={css}
			messages={messages}
			loading={loading}
			error={error}
			threadID={threadID}
			onToggleThreadClick={onToggleThreadClick}
			onEditClick={onEditClick}
			onDeleteClick={onDeleteClick}
			onCopyClick={onCopyClick}
		/>
	)
}

const mapStateToProps = state => ({
	threadID: selectThreadID(state),
})

const mapDispatchToProps = {
	setThreadMessage: setThreadMessageActionCreator,
	setEditMessage: setEditMessageActionCreator,
	deleteMessage: deleteMessageActionCreator,
	copyMessage: copyMessageActionCreator,
}

export default connect(
	mapStateToProps, mapDispatchToProps)(MessagesListContainer)