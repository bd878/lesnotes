import React, {Suspense, lazy, useState, useEffect, useCallback} from 'react';
import ReactDOM from 'react-dom/client';
import api from '../../api';
import Button from '../../components/Button';
import Tag from '../../components/Tag';
import Auth from '../../providers/Auth';
import i18n from '../../i18n';

const MessagesList = lazy(() => import("../../components/MessagesList"));
const SendMessageForm = lazy(() => import("../../components/SendMessageForm"));

const Messages = () => {
  const [error, setError] = useState(false)
  const [loading, setLoading] = useState(true)
  const [messages, setMessages] = useState([])

  const reload = useCallback(() => {
    const load = async () => {
      let response = []
      try {
        setLoading(true)
        response = await api("/messages/v1/read", {
          method: "GET",
          credentials: 'include',
        });
      } catch (e) {
        console.error(i18n("error_occured"), e);
        setError(i18n("loading_messages_error"));
      } finally {
        setLoading(false)
      }

      if (Array.isArray(response)) {
        setMessages(response);
      }
    }

    load();
  }, [setLoading, setMessages, setError])

  useEffect(reload, [reload]);

  const exit = useCallback(() => {
    /*TODO: implement*/
    setTimeout(() => {location.href = "/login"}, 0)
  }, []);

  const onSendSuccess = useCallback((response) => {
    console.log("[onSendSuccess] response:", response);
    setMessages([
      ...messages,
      response.message,
    ]);
  }, [setMessages, messages]);

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
          <MessagesList
            css="grow y-scroll hidden"
            error={error}
            messages={messages}
            setMessages={setMessages}
            loading={loading}
          />

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
