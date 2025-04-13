import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {
	setEditMessageActionCreator,
	deleteMessageActionCreator,
} from '../../features/messages';
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
	} = props

	const onOpenThreadClick = useCallback(setThreadMessage, [setThreadMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback((e) => {
		async function copy(e) {
			try {
				await navigator.clipboard.writeText("test")
			} catch (e) {
				console.error(e)
			}
		}

		copy(e)
	}, [])

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