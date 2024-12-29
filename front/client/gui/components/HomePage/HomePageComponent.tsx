import React, {lazy} from 'react';
import i18n from '../../i18n';
import Button from '../Button';
import Tag from '../Tag';

const MessagesList = lazy(() => import("../MessagesList"));
const SendMessageForm = lazy(() => import("../SendMessageForm"));

function HomePageComponent(props) {
  const {
    listRef,
    onExitClick,
    onListScroll,
    onLoadMoreClick,
    isAllLoaded,
    error,
    messages,
    loading,
    onSendSuccess,
    onSendError,
  } = props;

  return (
    <>
      <Button
        text={i18n("logout")}
        onClick={onExitClick}
      />

      <Tag css="flex column grow y-hidden w-100">
        <Tag>{i18n("messages_header")}</Tag>
        <Button
          text={i18n("load_more")}
          onClick={onLoadMoreClick}
          disabled={isAllLoaded}
        />

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
    </>
  )
}

export default HomePageComponent;
