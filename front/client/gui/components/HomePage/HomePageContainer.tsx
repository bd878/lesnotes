import React, {useRef, useEffect, useCallback} from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../third_party/react-redux';
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
} from '../../features/messages';
import {
	selectThreadID
} from '../../features/threads';
import {selectUser, logoutActionCreator} from '../../features/me';

function HomePageContainer(props) {
	const {
		messages,
		user,
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
	} = props

	const listRef = useRef(null);

	const scrollToTop = useCallback(() => {
		if (listRef.current != null) {
			listRef.current.scrollTo(0, listRef.current.scrollHeight);
		}
	}, [listRef]);

	useEffect(() => {
		fetchMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER)
	}, [fetchMessages]);

	const loadMore = useCallback(() => {
		if (listRef.current != null && !isLoading && !isLastPage) {
			fetchMessages(LIMIT_LOAD_BY, loadOffset, LOAD_ORDER)
		}
	}, [listRef.current, fetchMessages,
		loadOffset, isLoading, isLastPage]);

	const onListScroll = useCallback(() => {
		if (listRef.current != null && listRef.current.scrollTop == 0) {
			loadMore()
		}
	}, [listRef.current, loadMore]);

	const onExitClick = useCallback(() => {logout()}, [logout]);

	return (
		<HomePageComponent
			listRef={listRef}
			onExitClick={onExitClick}
			onListScroll={onListScroll}
			onLoadMoreClick={loadMore}
			isAllLoaded={isLastPage}
			error={error}
			messages={messages}
			loading={isLoading}
			sendMessage={sendMessage}
			updateMessage={updateMessage}
			resetEditMessage={resetEditMessage}
			messageForEdit={messageForEdit}
		/>
	)
}

const mapStateToProps = state => ({
	messages: selectMessages(state),
	isLoading: selectIsLoading(state),
	isLastPage: selectIsLastPage(state),
	loadOffset: selectLoadOffset(state),
	error: selectError(state),
	user: selectUser(state),
	messageForEdit: selectMessageForEdit(state),
	isEditMode: selectIsEditMode(state),
})

const mapDispatchToProps = {
	fetchMessages: fetchMessagesActionCreator,
	logout: logoutActionCreator,
	sendMessage: sendMessageActionCreator,
	updateMessage: updateMessageActionCreator,
	resetEditMessage: resetEditMessageActionCreator,
}

export default connect(
	mapStateToProps,
	mapDispatchToProps,
)(HomePageContainer);
