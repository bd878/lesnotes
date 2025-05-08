import React, {lazy, forwardRef} from 'react';
import Button from '../Button';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is';

const MessagesList = lazy(() => import("../MessagesList"));
const MessageElement = lazy(() => import("../MessageElement"))
const MessageForm = lazy(() => import("../MessageForm"));
const MainMessage = lazy(() => import("../MainMessage"));

function ThreadComponent(props, ref) {
	const {
		css,
		destroyContent,
		onDestroyClick,
		onDeleteSelectedClick,
		onLoadMoreClick,
		onSelectClick,
		onUnselectClick,
		onClearSelectedClick,
		onPublishClick,
		onPrivateClick,
		isAllLoaded,
		onScroll,
		loadMoreContent,
		error,
		loading,
		messages,
		selectedMessageIDs,
		isAnyMessageSelected,
		isAnyOpen,
		checkMyThreadOpen,
		onEditClick,
		onToggleThreadClick,
		onCopyClick,
		send,
		update,
		reset,
		messageForEdit,
	} = props

	return (
		<Tag css={(css || "") + " " + "flex flex-col items-start w-lg w-full"}>
			<Button
				content={destroyContent}
				onClick={onDestroyClick}
				css="btn mb-2"
			/>

			<Button
				tabIndex="0"
				content={loadMoreContent}
				css="disabled:opacity-30 btn w-full text-center mb-5"
				onClick={onLoadMoreClick}
				disabled={isAllLoaded}
			/>

			<Tag
				ref={ref}
				css="grow w-full h-full overflow-x-hidden overflow-y-scroll relative"
				onScroll={onScroll}
			>
				{isAnyMessageSelected ? (
					<Tag css="absolute w-full">
						<Button tabIndex="0" content={i18n("delete_messages")} onClick={onDeleteSelectedClick} />
						<Button tabIndex="0" content={i18n("cancel_delete")} onClick={onClearSelectedClick} />
					</Tag>
				) : null}

				<MessagesList
					css="w-full"
					error={error}
					messages={messages}
					messageForEdit={messageForEdit}
					selectedMessageIDs={selectedMessageIDs}
					loading={loading}
					isAnyThreadOpen={isAnyOpen}
					checkMyThreadOpen={checkMyThreadOpen}
					onPublishClick={onPublishClick}
					onPrivateClick={onPrivateClick}
					onSelectClick={onSelectClick}
					onUnselectClick={onUnselectClick}
					onEditClick={onEditClick}
					onResetEditClick={reset}
					onToggleThreadClick={onToggleThreadClick}
					onCopyClick={onCopyClick}
				/>
			</Tag>

			<Tag css="w-full mt-5">
				<MessageForm
					key={is.notEmpty(messageForEdit) ? messageForEdit.ID : null}
					send={send}
					update={update}
					reset={reset}
					messageForEdit={messageForEdit}
				/>
			</Tag>
		</Tag>
	)
}

export default forwardRef(ThreadComponent);
