import React, {useCallback} from 'react';
import MessagesListComponent from './MessagesListComponent';
import {setEditMessageActionCreator} from '../../features/messages';
import {connect} from '../../third_party/react-redux';

function MessagesListContainer(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		setEditMessage,
	} = props

	const onListItemClick = useCallback(setEditMessage, [setEditMessage])

	return (
		<MessagesListComponent
			css={css}
			liCss={liCss}
			messages={messages}
			loading={loading}
			error={error}
			onListItemClick={onListItemClick}
		/>
	)
}

const mapStateToProps = () => {}

const mapDispatchToProps = {
	setEditMessage: setEditMessageActionCreator,
}

export default connect(
	mapStateToProps, mapDispatchToProps)(MessagesListContainer)