import React, {lazy,useEffect,useState} from 'react';
import api from '../../api';

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
        console.error("error occured:", e);
        setError("error occured while loading messages");
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
      <div>Messages:</div>

      {loading ? <div>Loading...</div> : <>{content}</>}
    </>
  )
}

export default MessagesList;
