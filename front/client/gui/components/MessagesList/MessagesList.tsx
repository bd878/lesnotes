import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../i18n';
import {getFileDownloadUrl} from "../../api";

const List = lazy(() => import("../../components/List"));
const ListItem = lazy(() => import("../../components/ListItem"));

function MessagesList(props) {
  const {
    css,
    liCss,
    messages,
    loading,
    error,
  } = props

  let content = <Tag></Tag>;
  if (error) {
    content = <Tag>{error}</Tag>
  } else {
    content = (
      <List el="ul" css={css}>
        {messages.map(message => (
          <Tag
            el="li"
            css={liCss}
            key={message.id}
          >
            {(message.file && message.file.id && message.file.name) ? <Tag
              el="a"
              href={getFileDownloadUrl(`/files/v1/${message.file.id}`, false)}
              download={message.file.name}
              target="_blank"
            >
              {message.file.name}
            </Tag> : null}

            <ListItem key={message.id}>{message.text}</ListItem>
          </Tag>
        ))}
      </List>
    )
  }

  return (
    <>
      {loading ? <Tag>{i18n("loading")}</Tag> : null}
      {content}
    </>
  )
}

export default MessagesList;
