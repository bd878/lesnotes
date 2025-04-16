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
	selectHasNextThread,
	selectThreadID,
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
		index,
		threadID,
		openThread,
		closeThread,
		messages,
		error,
		logout,
		isLastPage,
		isLoading,
		hasNextThread,
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
		fetchMessages({limit: LIMIT_LOAD_BY, offset: 0, order: LOAD_ORDER, threadID: threadID})
	}, [fetchMessages, threadID]);

	const loadMore = useCallback(() => {
		if (is.notEmpty(listRef.current) && !isLoading && !isLastPage)
			fetchMessages({limit: LIMIT_LOAD_BY, offset: loadOffset, order: LOAD_ORDER})
	}, [listRef.current, fetchMessages, loadOffset, isLoading, isLastPage]);

	const onListScroll = useCallback(() => {
		if (is.notEmpty(listRef.current) && is.notEmpty(listRef.current.scrollTop))
			loadMore()
	}, [listRef.current, loadMore]);

	const onToggleThreadClick = useCallback(message => {
		if (hasNextThread)
			closeThread(message)
		else
			openThread(message)
	}, [closeThread, openThread, hasNextThread])

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
			onToggleThreadClick={onToggleThreadClick}
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
	hasNextThread: selectHasNextThread(index)(state),
	isLoading: selectIsLoading(index)(state),
	isLastPage: selectIsLastPage(index)(state),
	loadOffset: selectLoadOffset(index)(state),
	error: selectError(index)(state),
	messageForEdit: selectMessageForEdit(index)(state),
	isEditMode: selectIsEditMode(index)(state),
	threadID: selectThreadID(index)(state),
})

const mapDispatchToProps = (dispatch, {index}) => ({
	fetchMessages: payload => dispatch(fetchMessagesActionCreator(index)(payload)),
	sendMessage: payload => dispatch(sendMessageActionCreator(index)(payload)),
	updateMessage: payload => dispatch(updateMessageActionCreator(index)(payload)),
	resetEditMessage: payload => dispatch(resetEditMessageActionCreator(index)(payload)),
	deleteMessage: payload => dispatch(deleteMessageActionCreator(index)(payload)),
	copyMessage: payload => dispatch(copyMessageActionCreator(index)(payload)),
	setEditMessage: payload => dispatch(setEditMessageActionCreator(index)(payload)),
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadContainer)