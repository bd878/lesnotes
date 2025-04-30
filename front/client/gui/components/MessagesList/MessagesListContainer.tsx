import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {
	setEditMessageActionCreator,
	deleteMessageActionCreator,
	copyMessageActionCreator,
	selectMessageActionCreator,
} from '../../features/stack';
import {
	setThreadMessageActionCreator,
} from '../../features/threads';
import {connect} from '../../../third_party/react-redux';

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
		copyMessage,
		selectMessage
		checkMyThreadOpen,
	} = props

	const onOpenThreadClick = useCallback(setThreadMessage, [setThreadMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback(copyMessage, [copyMessage])
	const onSelectClick = useCallback(selectMessage, [selectMessage])

	return (
		<MessagesListComponent
			css={css}
			liCss={liCss}
			messages={messages}
			loading={loading}
			error={error}
			onOpenThreadClick={onOpenThreadClick}
			onEditClick={onEditClick}
			onDeleteClick={onDeleteClick}
			onCopyClick={onCopyClick}
			checkMyThreadOpen={checkMyThreadOpen}
		/>
	)
}

const mapStateToProps = () => {}

const mapDispatchToProps = {
	setThreadMessage: setThreadMessageActionCreator,
	setEditMessage: setEditMessageActionCreator,
	deleteMessage: deleteMessageActionCreator,
	copyMessage: copyMessageActionCreator,
	selectMessage: selectMessageActionCreator,
}

export default connect(
	mapStateToProps, mapDispatchToProps)(MessagesListContainer)