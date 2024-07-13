import React, {lazy} from 'react';
import i18n from '../../i18n';

const List = lazy(() => import("../../components/List/List.tsx"));
const ListItem = lazy(() => import("../../components/ListItem/ListItem.tsx"));

const MessagesList = ({
  messages,
  setMessages,
  loading,
  error,
}) => {
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
