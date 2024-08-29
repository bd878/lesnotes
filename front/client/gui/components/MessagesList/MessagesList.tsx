import React, {lazy, useRef, useEffect} from 'react';
import Tag from '../Tag';
import i18n from '../../i18n';
import {getFileDownloadUrl} from "../../api";

const List = lazy(() => import("../../components/List"));
const ListItem = lazy(() => import("../../components/ListItem"));

const MessagesList = ({
  css,
  messages,
  setMessages,
  loading,
  error,
}) => {
  const listRef = useRef(null);

  useEffect(() => {
    if (listRef.current != null) {
      listRef.current.scrollTo(0, listRef.current.scrollHeight);
    }
  }, [
    messages.length,
    listRef.current,
  ]);

  let content = <Tag></Tag>;
  if (error) {
    content = <Tag>{error}</Tag>
  } else {
    content = (
      <List ref={listRef} el="ul" css={css}>
        {messages.map(message => (
          <Tag el="li" key={message.id}>
            <Tag
              el="a"
              href={getFileDownloadUrl(`/messages/v1/read_file?id=${message.fileid}`, false)}
              download={message.filename}
              target="_blank"
            >
              {message.filename}
            </Tag>

            <ListItem key={message.id}>{message.value}</ListItem>
          </Tag>
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
