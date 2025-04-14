import React, {useRef, useEffect, useCallback} from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../../third_party/react-redux';
import * as is from '../../../third_party/is';
import {equal} from '../../../utils';
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
} from '../../features/messages';
import {
	fetchMessagesActionCreator as fetchThreadMessagesActionCreator,
	sendMessageActionCreator as sendThreadMessageActionCreator,
	resetActionCreator as resetThreadActionCreator,
	selectMessages as selectThreadMessages,
	selectIsLoading as selectThreadLoading,
	selectError as selectThreadError,
	selectThreadMessage,
	selectThreadID,
	setThreadMessageActionCreator,
} from '../../features/threads';
import {selectUser, logoutActionCreator} from '../../features/me';

function HomePageContainer(props) {
	const {
		messages,
		user,
		error,
		threadError,
		logout,
		isLastPage,
		isLoading,
		isThreadLoading,
		loadOffset,
		fetchMessages,
		sendMessage,
		updateMessage,
		resetEditMessage,
		messageForEdit,
		setThreadMessage,
		setEditMessage,
		deleteMessage,
		copyMessage,

		threadMessage,
		threadMessages,
		threadID,
		resetThread,
		fetchThreadMessages,
		sendThreadMessage,
	} = props

	const listRef = useRef(null);
	const threadListRef = useRef(null);

	const scrollToTop = useCallback(() => {
		if (is.notEmpty(listRef.current))
			listRef.current.scrollTo(0, listRef.current.scrollHeight);
	}, [listRef]);

	useEffect(() => {
		fetchMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [fetchMessages]);

	useEffect(() => {
		if (is.notEmpty(threadID))
			fetchThreadMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [threadID])

	const loadMore = useCallback(() => {
		if (is.notEmpty(listRef.current) && !isLoading && !isLastPage)
			fetchMessages(LIMIT_LOAD_BY, loadOffset, LOAD_ORDER)
	}, [listRef.current, fetchMessages,
		loadOffset, isLoading, isLastPage]);

	const onListScroll = useCallback(() => {
		if (is.notEmpty(listRef.current) && is.notEmpty(listRef.current.scrollTop))
			loadMore()
	}, [listRef.current, loadMore]);

	const onThreadListScroll = useCallback(() => {}, [])

	const onThreadLoadMore = useCallback(() => {

	}, [])

	const onExitClick = useCallback(logout, [logout]);

	const onCloseThreadClick = useCallback(resetThread, [resetThread])

	const onToggleThreadClick = useCallback(setThreadMessage, [setThreadMessage])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback(copyMessage, [copyMessage])

	return (
		<HomePageComponent
			listRef={listRef}
			threadListRef={threadListRef}
			onExitClick={onExitClick}
			onListScroll={onListScroll}
			onThreadListScroll={onThreadListScroll}
			onLoadMoreClick={loadMore}
			isAllLoaded={isLastPage}
			error={error}
			messages={messages}
			loading={isLoading}
			threadLoading={isThreadLoading}
			threadError={threadError}
			sendMessage={sendMessage}
			updateMessage={updateMessage}
			resetEditMessage={resetEditMessage}
			messageForEdit={messageForEdit}
			shouldShowThreadsPanel={threadID != 0}
			threadID={threadID}
			checkMessageThreadOpen={equal(threadID)}
			onLoadMoreThreadMessagesClick={() => {}}
			isAllThreadMessagesLoaded={false}
			threadMessage={threadMessage}
			threadMessages={threadMessages}
			closeThread={onCloseThreadClick}
			sendThreadMessage={sendThreadMessage}
			onToggleThreadClick={onToggleThreadClick}
			onDeleteClick={onDeleteClick}
			onEditClick={onEditClick}
			onCopyClick={onCopyClick}
		/>
	)
}

const mapStateToProps = state => ({
	messages: selectMessages(state),
	isLoading: selectIsLoading(state),
	isLastPage: selectIsLastPage(state),
	isThreadLoading: selectThreadLoading(state),
	threadError: selectThreadError(state),
	loadOffset: selectLoadOffset(state),
	error: selectError(state),
	user: selectUser(state),
	messageForEdit: selectMessageForEdit(state),
	isEditMode: selectIsEditMode(state),

	threadMessage: selectThreadMessage(state),
	threadMessages: selectThreadMessages(state),
	threadID: selectThreadID(state),
})

const mapDispatchToProps = {
	fetchMessages: fetchMessagesActionCreator,
	logout: logoutActionCreator,
	sendMessage: sendMessageActionCreator,
	updateMessage: updateMessageActionCreator,
	resetEditMessage: resetEditMessageActionCreator,
	deleteMessage: deleteMessageActionCreator,
	copyMessage: copyMessageActionCreator,
	setEditMessage: setEditMessageActionCreator,
	setThreadMessage: setThreadMessageActionCreator,

	resetThread: resetThreadActionCreator,
	fetchThreadMessages: fetchThreadMessagesActionCreator,
	sendThreadMessage: sendThreadMessageActionCreator,
}

export default connect(
	mapStateToProps,
	mapDispatchToProps,
)(HomePageContainer);
