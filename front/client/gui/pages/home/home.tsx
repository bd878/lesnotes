import React, {Suspense, lazy, useState, useRef, useEffect, useCallback} from 'react';
import ReactDOM from 'react-dom/client';
import api from '../../api';
import Button from '../../components/Button';
import Tag from '../../components/Tag';
import Auth from '../../providers/Auth';
import i18n from '../../i18n';
import throttle from '../../utils/throttle';

const MessagesList = lazy(() => import("../../components/MessagesList"));
const SendMessageForm = lazy(() => import("../../components/SendMessageForm"));

const TOP_GAP_TO_LOAD_MORE = 70;
const LIMIT_LOAD_BY = 25;
const LOAD_ORDER = 0; /* descending */

async function loadMessages(limit, offset) {
  let response = [];
  try {
    response = await api('/messages/v1/read', {
      queryParams: {
        limit: limit,
        offset: offset,
        asc: LOAD_ORDER,
      },
      method: "GET",
      credentials: 'include',
    });
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  if (Array.isArray(response.messages)) {
    return response
  }

  return {
    messages: [],
    islastpage: true,
  };
}

const Messages = () => {
  const listRef = useRef(null);
  const [error, setError] = useState(false);
  const [isLastPage, setIsLastPage] = useState(false);
  const [loading, setLoading] = useState(true);
  const [messages, setMessages] = useState([]);

  const getLoadOffset = useCallback(() => messages.length, [messages]);

  const scrollToTop = useCallback(() => {
    if (listRef.current != null) {
      listRef.current.scrollTo(0, listRef.current.scrollHeight);
    }
  }, [listRef]);

  const appendMessages = useCallback((messagesToAppend) => {
    setMessages([ ...messages, ...messagesToAppend ]);
  }, [messages, setMessages]);

  const pushBackMessages = useCallback((messagesToPushBack) => {
    setMessages([ ...messagesToPushBack, ...messages ]);
  }, [messages, setMessages]);

  useEffect(() => {
    const init = async () => {
      try {
        setLoading(true);
        const response = await loadMessages(LIMIT_LOAD_BY, 0);
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

  const onListScroll = useCallback(() => {
    const loadMore = async () => {
      try {
        setLoading(true);
        console.log("load more")
        const response = await loadMessages(LIMIT_LOAD_BY, getLoadOffset())
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
      !isLastPage &&
      (listRef.current.scrollTop == 0)
    ) {
      loadMore()
    }
  }, [
    listRef.current,
    loading,
    isLastPage,
    getLoadOffset,
    setIsLastPage,
  ]);

  const exit = useCallback(() => {
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
    <Auth fallback={i18n("messages_auth_fallback")}>
      <Suspense fallback={i18n("loading")}>
        <Button
          text={i18n("logout")}
          onClick={exit}
        />

        <Tag css="flex column grow y-hidden w-100">
          <Tag>{i18n("messages_header")}</Tag>

          <Tag
            el="div"
            ref={listRef}
            css="grow y-scroll"
            onScroll={onListScroll}
          >
            <MessagesList
              css="reset-list-style"
              liCss="li-10"
              error={error}
              messages={messages}
              loading={loading}
            />
          </Tag>

          <SendMessageForm
            onSuccess={onSendSuccess}
            onError={onSendError}
          />
        </Tag>
      </Suspense>
    </Auth>
  )
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
