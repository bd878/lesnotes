import React, {lazy,useEffect,useState} from 'react';
import api from '../../api';
import i18n from '../../i18n';

const List = lazy(() => import("../../components/List/index.jsx"));
const ListItem = lazy(() => import("../../components/ListItem/index.jsx"));

const MessagesList = () => {
  const [error, setError] = useState(false)
  const [loading, setLoading] = useState(true)
  const [messages, setMessages] = useState([])

  useEffect(() => {
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
  }, [setLoading, setMessages, setError]);

  let content = <div></div>;
  if (error) {
    content = <div>{error}</div>
  } else {
    content = (
      <List>
        {messages.map(message => (
          <ListItem key={message.id}>{message.value}</ListItem>
        ))}
      </List>
    )
  }

  return (
    <>
      <div>{i18n("messages_header")}</div>

      {loading ? <div>{i18n("loading")}</div> : <>{content}</>}
    </>
  )
}

export default MessagesList;
