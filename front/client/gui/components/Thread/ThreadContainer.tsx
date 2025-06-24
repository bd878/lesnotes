import React, {useEffect, useCallback, useMemo, useRef} from 'react';
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
	selectIsMessageThreadOpen,
	selectSelectedMessageIDs,
	sendMessageActionCreator,
	updateMessageActionCreator,
	selectMessageForEdit,
	selectIsEditMode,
	resetEditMessageActionCreator,
	setEditMessageActionCreator,
	deleteMessageActionCreator,
	deleteSelectedActionCreator,
	copyMessageActionCreator,
	selectMessageActionCreator,
	unselectMessageActionCreator,
	clearSelectedActionCreator,
	publishSelectedActionCreator,
	privateSelectedActionCreator,
	copyLinkActionCreator,
	moveMessageActionCreator,
} from '../../features/stack';

function ThreadContainer(props) {
	const {
		css,
		index,
		destroyContent,
		checkMyThreadOpen,
		threadID,
		openThread,
		closeThread,
		copyLink,
		destroyThread,
		messages,
		selectedMessageIDs,
		error,
		logout,
		isLastPage,
		isLoading,
		hasNextThread,
		loadOffset,
		fetchMessages,
		sendMessage,
		updateMessage,
		moveMessage,
		resetEditMessage,
		messageForEdit,
		setEditMessage,
		deleteMessage,
		deleteSelected,
		copyMessage,
		selectMessage,
		unselectMessage,
		clearSelected,
		publishSelected,
		privateSelected,
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
		if (!isLoading && !isLastPage)
			fetchMessages({limit: LIMIT_LOAD_BY, offset: loadOffset, order: LOAD_ORDER, threadID: threadID})
	}, [listRef.current, fetchMessages, loadOffset, isLoading, isLastPage, threadID]);

	const onListScroll = useCallback(() => {
		if (is.notEmpty(listRef.current) && listRef.current.scrollTop == 0)
			loadMore()
	}, [listRef.current, loadMore]);

	const onToggleThreadClick = useCallback(message => {
		if (hasNextThread)
			closeThread(message)
		else
			openThread(message)
	}, [closeThread, openThread, hasNextThread])

	const onDrop = useCallback((event) => {
		event.stopPropagation()
		const [prevIndex, message] = JSON.parse(event.dataTransfer.getData('text/plain'))

		moveMessage(prevIndex, message)

		return false
	}, [moveMessage])

	const onDragOver = useCallback(event => {
		const [prevIndex, message] = JSON.parse(event.dataTransfer.getData('text/plain'))

		if (prevIndex == index)
			return false

		if (threadID === message.ID)
			return false

		event.preventDefault()

		return false
	}, [index, threadID])

	const onDeleteSelectedClick = useCallback(deleteSelected, [deleteSelected])
	const onDeleteClick = useCallback(deleteMessage, [deleteMessage])
	const onEditClick = useCallback(setEditMessage, [setEditMessage])
	const onCopyClick = useCallback(copyMessage, [copyMessage])
	const onSelectClick = useCallback(selectMessage, [selectMessage])
	const onUnselectClick = useCallback(unselectMessage, [unselectMessage])
	const onClearSelectedClick = useCallback(clearSelected, [clearSelected])
	const onPublishClick = useCallback(publishSelected, [publishSelected])
	const onPrivateClick = useCallback(privateSelected, [privateSelected])
	const onCopyLinkClick = useCallback(copyLink, [copyLink])

	const onMessageSend = useCallback(payload => {
		payload.threadID = threadID
		sendMessage(payload)
	}, [sendMessage, threadID])

	const isAnyMessageSelected = useMemo(() => selectedMessageIDs.size > 0, [selectedMessageIDs])

	return (
		<ThreadComponent
			ref={listRef}
			css={css}
			index={index}
			destroyContent={destroyContent}
			loadMoreContent={i18n("load_more")}
			onDestroyClick={destroyThread}
			onLoadMoreClick={loadMore}
			isAllLoaded={isLastPage}
			onScroll={onListScroll}
			onDrop={onDrop}
			onDragOver={onDragOver}
			error={error}
			loading={isLoading}
			isAnyMessageSelected={isAnyMessageSelected}
			selectedMessageIDs={selectedMessageIDs}
			messages={messages}
			isAnyOpen={hasNextThread}
			checkMyThreadOpen={checkMyThreadOpen}
			onSelectClick={onSelectClick}
			onUnselectClick={onUnselectClick}
			onClearSelectedClick={onClearSelectedClick}
			onDeleteSelectedClick={onDeleteSelectedClick}
			onDeleteClick={onDeleteClick}
			onCopyLinkClick={onCopyLinkClick}
			onEditClick={onEditClick}
			onPublishClick={onPublishClick}
			onPrivateClick={onPrivateClick}
			onToggleThreadClick={onToggleThreadClick}
			onCopyClick={onCopyClick}
			send={onMessageSend}
			update={updateMessage}
			reset={resetEditMessage}
			messageForEdit={messageForEdit}
		/>
	)
}

const mapStateToProps = (state, {index}) => ({
	messages: selectMessages(index)(state),
	selectedMessageIDs: selectSelectedMessageIDs(index)(state),
	hasNextThread: selectHasNextThread(index)(state),
	isLoading: selectIsLoading(index)(state),
	isLastPage: selectIsLastPage(index)(state),
	loadOffset: selectLoadOffset(index)(state),
	error: selectError(index)(state),
	messageForEdit: selectMessageForEdit(index)(state),
	isEditMode: selectIsEditMode(index)(state),
	threadID: selectThreadID(index)(state),
	checkMyThreadOpen: (messageID) => selectIsMessageThreadOpen(index)(state)(messageID),
})

const mapDispatchToProps = (dispatch, {index}) => ({
	moveMessage: (prevIndex, payload) => dispatch(moveMessageActionCreator(index)(prevIndex, payload)),
	publishSelected: payload => dispatch(publishSelectedActionCreator(index)(payload)),
	privateSelected: payload => dispatch(privateSelectedActionCreator(index)(payload)),
	clearSelected: payload => dispatch(clearSelectedActionCreator(index)(payload)),
	unselectMessage: payload => dispatch(unselectMessageActionCreator(index)(payload)),
	selectMessage: payload => dispatch(selectMessageActionCreator(index)(payload)),
	fetchMessages: payload => dispatch(fetchMessagesActionCreator(index)(payload)),
	sendMessage: payload => dispatch(sendMessageActionCreator(index)(payload)),
	updateMessage: payload => dispatch(updateMessageActionCreator(index)(payload)),
	resetEditMessage: payload => dispatch(resetEditMessageActionCreator(index)(payload)),
	deleteSelected: payload => dispatch(deleteSelectedActionCreator(index)(payload)),
	deleteMessage: payload => dispatch(deleteMessageActionCreator(index)(payload)),
	copyMessage: payload => dispatch(copyMessageActionCreator(index)(payload)),
	copyLink: payload => dispatch(copyLinkActionCreator(index)(payload)),
	setEditMessage: payload => dispatch(setEditMessageActionCreator(index)(payload)),
})

export default connect(mapStateToProps, mapDispatchToProps)(ThreadContainer)