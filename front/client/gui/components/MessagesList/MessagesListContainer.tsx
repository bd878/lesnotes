import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {
	setEditMessageActionCreator,
	deleteMessageActionCreator,
} from '../../features/messages';
import {
	setThreadMessageActionCreator,
} from '../../features/threads';
import {connect} from '../../third_party/react-redux';

function MessagesListContainer(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		setThreadMessage,
		setEditMessage,
		deleteMessage,
	} = props

	const onListItemClick = useCallback(setThreadMessage, [setThreadMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])

	return (
		<MessagesListComponent
			css={css}
			liCss={liCss}
			messages={messages}
			loading={loading}
			error={error}
			onListItemClick={onListItemClick}
			onEditClick={onEditClick}
			onDeleteClick={onDeleteClick}
		/>
	)
}

const mapStateToProps = () => {}

const mapDispatchToProps = {
	setThreadMessage: setThreadMessageActionCreator,
	setEditMessage: setEditMessageActionCreator,
	deleteMessage: deleteMessageActionCreator,
}

export default connect(
	mapStateToProps, mapDispatchToProps)(MessagesListContainer)