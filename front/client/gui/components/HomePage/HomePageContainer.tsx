import React, {useState, useRef, useEffect, useCallback} from 'react';
import i18n from '../../i18n';
import HomePageComponent from './HomePageComponent';
import loadMessages from './loadMessages';
import {
  LIMIT_LOAD_BY,
  LOAD_ORDER,
} from './const';

function HomePageContainer(props) {
  const {
    messages,
    appendMessages,
    pushBackMessages,
  } = props

  const listRef = useRef(null);
  const [error, setError] = useState(false);
  const [isLastPage, setIsLastPage] = useState(false);
  const [loading, setLoading] = useState(true);

  const getLoadOffset = useCallback(() => messages.length, [messages]);

  const scrollToTop = useCallback(() => {
    if (listRef.current != null) {
      listRef.current.scrollTo(0, listRef.current.scrollHeight);
    }
  }, [listRef]);

  useEffect(() => {
    const init = async () => {
      try {
        setLoading(true);
        const response = await loadMessages(LIMIT_LOAD_BY, 0, LOAD_ORDER);
        if (response.error != "") {
          console.error("[HomePageContainer]: failed to load messages",
            response.error, response.explain);
          throw(response.error);
        }
        response.messages.reverse();
        appendMessages(response.messages);
        if (response.islastpage) {
          setIsLastPage(true)
        }
        setTimeout(scrollToTop, 300);
      } catch (_1) {
        setError(true)
      } finally {
        setLoading(false);
      }
    }

    init();
  }, []);

  const loadMore = useCallback(() => {
    const load = async () => {
      try {
        setLoading(true);
        const response = await loadMessages(
          LIMIT_LOAD_BY,
          getLoadOffset(),
          LOAD_ORDER,
        )
        if (response.error != "") {
          console.error("[HomePageContainer]: failed to load messages",
            response.error, response.explain);
          throw(response.error);
        }
        response.messages.reverse();
        if (response.islastpage) {
          setIsLastPage(true)
        }
        pushBackMessages(response.messages);
      } catch (_1) {
        setError(true)
      } finally {
        setLoading(false);
      }
    }

    if (
      (listRef.current != null) && 
      !loading &&
      !isLastPage
    ) {
      load()
    }
  }, [
    listRef.current,
    loading,
    isLastPage,
    getLoadOffset,
    setIsLastPage,
  ]);

  const onListScroll = useCallback(() => {
    if (
      listRef.current != null &&
      (listRef.current.scrollTop == 0)
    ) {
      loadMore()
    }
  }, [listRef.current, loadMore]);

  const onExitClick = useCallback(() => {
    /*TODO: implement*/
    setTimeout(() => {location.href = "/login"}, 0)
  }, []);

  const onSendSuccess = useCallback((newMessage) => {
    appendMessages([newMessage]);
    setTimeout(scrollToTop, 0);
  }, [appendMessages]);

  const onSendError = useCallback(() => {
    setError(i18n("loading_messages_error"))
  }, [setError])

  return (
    <HomePageComponent
      listRef={listRef}
      onExitClick={onExitClick}
      onSendSuccess={onSendSuccess}
      onSendError={onSendError}
      onListScroll={onListScroll}
      onLoadMoreClick={loadMore}
      isAllLoaded={isLastPage}
      error={error}
      messages={messages}
      loading={loading}
    />
  )
}

export default HomePageContainer;
