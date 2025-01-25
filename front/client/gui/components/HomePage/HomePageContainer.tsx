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
} from '../../features/messages';

function HomePageContainer(props) {
  const {
    messages,
    error,
    isLastPage,
    isLoading,
    loadOffset,
    fetchMessages,
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

  const onExitClick = useCallback(() => {
    /*TODO: implement*/
    setTimeout(() => {location.href = "/login"}, 0)
  }, []);

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
    />
  )
}

const mapStateToProps = state => ({
  messages: selectMessages(state),
  isLoading: selectIsLoading(state),
  isLastPage: selectIsLastPage(state),
  loadOffset: selectLoadOffset(state),
  error: selectError(state),
})

const mapDispatchToProps = {
  fetchMessages: fetchMessagesActionCreator,
}

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(HomePageContainer);
