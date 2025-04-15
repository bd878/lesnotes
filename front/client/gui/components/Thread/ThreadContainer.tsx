import React, {useEffect, useCallback, useRef} from 'react';
import {connect} from '../../../third_party/react-redux';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is';
import ThreadComponent from './ThreadComponent'
import {
	LIMIT_LOAD_BY,
	LOAD_ORDER,
} from './const';
import {
	fetchMessagesActionCreator,
	selectMessages,
	selectIsLastPage,
	selectIsLoading,
	selectError,
	selectLoadOffset,
	sendMessageActionCreator,
	updateMessageActionCreator,
	selectMessageForEdit,
	selectIsEditMode,
	resetEditMessageActionCreator,
	setEditMessageActionCreator,
	deleteMessageActionCreator,
	copyMessageActionCreator,
} from '../../features/stack';

function ThreadContainer(props) {
	const {
		messages,
		error,
		logout,
		isLastPage,
		isLoading,
		loadOffset,
		fetchMessages,
		sendMessage,
		updateMessage,
		resetEditMessage,
		messageForEdit,
		setEditMessage,
		deleteMessage,
		copyMessage,
	} = props

	const listRef = useRef(null);

	const scrollToTop = useCallback(() => {
		if (is.notEmpty(listRef.current))
			listRef.current.scrollTo(0, listRef.current.scrollHeight);
	}, [listRef]);

	useEffect(() => {
		fetchMessages({limit: LIMIT_LOAD_BY, offset: 0, order: LOAD_ORDER})
	}, [fetchMessages]);

	const loadMore = useCallback(() => {
		if (is.notEmpty(listRef.current) && !isLoading && !isLastPage)
			fetchMessages({limit: LIMIT_LOAD_BY, offset: loadOffset, order: LOAD_ORDER})
	}, [listRef.current, fetchMessages, loadOffset, isLoading, isLastPage]);

	const onListScroll = useCallback(() => {
		if (is.notEmpty(listRef.current) && is.notEmpty(listRef.current.scrollTop))
			loadMore()
	}, [listRef.current, loadMore]);

	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback(copyMessage, [copyMessage])

	return (
		<ThreadComponent
			ref={listRef}
			destroyContent={"< " + i18n("logout")}
			loadMoreContent={i18n("load_more")}
			onDestroyClick={() => {}}
			onLoadMoreClick={() => {}}
			isAllLoaded={false}
			onScroll={onListScroll}
			error={error}
			loading={isLoading}
			messages={messages}
			isAnyOpen={false}
			checkMyThreadOpen={() => {}}
			onDeleteClick={onDeleteClick}
			onEditClick={onEditClick}
			onToggleThreadClick={() => {}}
			onCopyClick={onCopyClick}
			send={sendMessage}
			update={updateMessage}
			reset={resetEditMessage}
			messageForEdit={messageForEdit}
		/>
	)
}

const mapStateToProps = (state, {index}) => ({
	messages: selectMessages(index)(state),
	isLoading: selectIsLoading(index)(state),
	isLastPage: selectIsLastPage(index)(state),
	loadOffset: selectLoadOffset(index)(state),
	error: selectError(index)(state),
	messageForEdit: selectMessageForEdit(index)(state),
	isEditMode: selectIsEditMode(index)(state),
})

const mapDispatchToProps = (_, {index}) => ({
	fetchMessages: fetchMessagesActionCreator(index),
	sendMessage: sendMessageActionCreator(index),
	updateMessage: updateMessageActionCreator(index),
	resetEditMessage: resetEditMessageActionCreator(index),
	deleteMessage: deleteMessageActionCreator(index),
	copyMessage: copyMessageActionCreator(index),
	setEditMessage: setEditMessageActionCreator(index),
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadContainer)