import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../i18n';

const List = lazy(() => import("../../components/List"));
const ListItem = lazy(() => import("../../components/ListItem"));

const MessagesList = ({
  css,
  messages,
  setMessages,
  loading,
  error,
}) => {
  let content = <Tag></Tag>;
  if (error) {
    content = <Tag>{error}</Tag>
  } else {
    content = (
      <List css={css}>
        {messages.map(message => (
          <ListItem key={message.id}>{message.value}</ListItem>
        ))}
      </List>
    )
  }

  return (
    <>
      <Tag>{i18n("messages_header")}</Tag>

      {loading ? <Tag>{i18n("loading")}</Tag> : <>{content}</>}
    </>
  )
}

export default MessagesList;
