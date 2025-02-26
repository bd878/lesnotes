import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {setEditMessageActionCreator, deleteMessageActionCreator} from '../../features/messages';
import {connect} from '../../third_party/react-redux';

function MessagesListContainer(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		setEditMessage,
		deleteMessage,
	} = props

	const onListItemClick = useCallback(setEditMessage, [setEditMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])

	return (
		<MessagesListComponent
			css={css}
			liCss={liCss}
			messages={messages}
			loading={loading}
			error={error}
			onListItemClick={onListItemClick}
			onDeleteClick={onDeleteClick}
		/>
	)
}

const mapStateToProps = () => {}

const mapDispatchToProps = {
	setEditMessage: setEditMessageActionCreator,
	deleteMessage: deleteMessageActionCreator,
}

export default connect(
	mapStateToProps, mapDispatchToProps)(MessagesListContainer)