import React, {Suspense, lazy, useState, useEffect, useCallback} from 'react';
import ReactDOM from 'react-dom/client';
import api from '../../api';
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

  return (
    <Auth fallback={i18n("messages_auth_fallback")}>
      <Suspense fallback={i18n("loading")}>
        <div>
          <MessagesList
            error={error}
            messages={messages}
            setMessages={setMessages}
            loading={loading}
          />

          <SendMessageForm
            onSend={reload}
            setError={setError}
          />
        </div>
      </Suspense>
    </Auth>
  )
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
