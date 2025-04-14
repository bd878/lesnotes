import React, {lazy} from 'react';
import i18n from '../../../i18n';
import Tag from '../Tag';
import * as is from '../../../third_party/is';

const Thread = lazy(() => import("../Thread"));

function HomePageComponent(props) {
	const {
		listRef,
		threadListRef,
		onExitClick,
		onListScroll,
		onThreadListScroll,
		onLoadMoreClick,
		isAllLoaded,
		error,
		threadError,
		isThreadLoading,
		messages,
		loading,
		sendMessage,
		updateMessage,
		resetEditMessage,
		messageForEdit,
		threadMessageForEdit,
		checkMessageThreadOpen,
		shouldShowThreadsPanel,
		threadID,
		threadMessage,
		threadMessages,
		onLoadMoreThreadMessagesClick,
		isAllThreadMessagesLoaded,
		closeThread,
		sendThreadMessage,
		onToggleThreadClick,
		onEditClick,
		onCopyClick,
		onDeleteClick,
		updateThreadMessage,
		resetEditThreadMessage,
	} = props;

	return (
		<>
			<Tag css="flex flex-row grow max-h-full pb-8">
				<Thread
					ref={listRef}
					destroyContent={"< " + i18n("logout")}
					loadMoreContent={i18n("load_more")}
					onDestroyClick={onExitClick}
					onLoadMoreClick={onLoadMoreClick}
					isAllLoaded={isAllLoaded}
					onScroll={onListScroll}
					error={error}
					loading={loading}
					messages={messages}
					isAnyOpen={is.notEmpty(threadID)}
					checkMyThreadOpen={checkMessageThreadOpen}
					onDeleteClick={onDeleteClick}
					onEditClick={onEditClick}
					onToggleThreadClick={onToggleThreadClick}
					onCopyClick={onCopyClick}
					send={sendMessage}
					update={updateMessage}
					reset={resetEditMessage}
					messageForEdit={messageForEdit}
				/>
			</Tag>
		</>
	)
}

export default HomePageComponent;
