import React, {lazy, forwardRef} from 'react';
import Button from '../Button';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is';
import RubbishIcon from '../../icons/RubbishIcon'
import EyeCloseIcon from '../../icons/EyeCloseIcon'
import EyeOpenIcon from '../../icons/EyeOpenIcon'
import BigCrossIcon from '../../icons/BigCrossIcon'

const MessagesList = lazy(() => import("../MessagesList"));
const MessageElement = lazy(() => import("../MessageElement"))
const MessageForm = lazy(() => import("../MessageForm"));
const MainMessage = lazy(() => import("../MainMessage"));

function ThreadComponent(props, ref) {
	const {
		css,
		index,
		destroyContent,
		onDestroyClick,
		onDeleteSelectedClick,
		onLoadMoreClick,
		onSelectClick,
		onUnselectClick,
		onClearSelectedClick,
		onPublishClick,
		onPrivateClick,
		onCopyLinkClick,
		isAllLoaded,
		onScroll,
		onDrop,
		onDragOver,
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
		<Tag css={(css || "") + " " + "flex flex-col items-start w-xl w-full relative"}>
			<Button
				content={destroyContent}
				onClick={onDestroyClick}
				css="btn mb-2"
			/>

			<Tag css="w-full shadow-[0_10px_10px_-10px_rgba(0,0,0,0.4)] z-1">
				<Button
					tabIndex="0"
					content={loadMoreContent}
					css="disabled:opacity-30 btn w-full text-center"
					onClick={onLoadMoreClick}
					disabled={isAllLoaded}
				/>
			</Tag>

			{isAnyMessageSelected ? (
				<Tag css="absolute top-[60px] w-full z-1 bg-white py-3 px-2 flex justify-between items-center shadow-[0_10px_10px_-10px_rgba(0,0,0,0.4)]">
					<Tag>
						<Button
							css="text-center cursor-pointer p-1 hover:text-red-800"
							tabIndex="0"
							onClick={onDeleteSelectedClick}
							content={
								<Tag css="flex flex-row items-center">
									<RubbishIcon css="mr-1" width="16" height="16" />
									<Tag el="span">{i18n("delete_message")}</Tag>
								</Tag>
							}
						/>
						{" / "}
						<Button
							css="text-center cursor-pointer p-1 hover:text-blue-800"
							onClick={onPrivateClick}
							content={
								<Tag css="flex flex-row items-center">
									<EyeCloseIcon css="mr-1" width="16" height="16" />
									<Tag el="span">{i18n("private_message")}</Tag>
								</Tag>
							}
						/>
						{" / "}
						<Button
							css="text-center cursor-pointer p-1 hover:text-blue-800"
							onClick={onPublishClick}
							content={
								<Tag css="flex flex-row items-center">
									<EyeOpenIcon css="mr-1" width="16" height="16" />
									<Tag el="span">{i18n("publish_message")}</Tag>
								</Tag>
							}
						/>
					</Tag>
					<Button
						css="p-2 rounded-sm cursor-pointer hover:bg-gray-300"
						tabIndex="0"
						content={<BigCrossIcon width="20" height="20" />}
						onClick={onClearSelectedClick}
					/>
				</Tag>
			) : null}

			<Tag
				ref={ref}
				css={"grow w-full flex flex-col overflow-x-hidden"}
				onScroll={onScroll}
				onDrop={onDrop}
				onDragOver={onDragOver}
			>
				<MessagesList
					css="w-full mt-[60px] mb-5"
					error={error}
					index={index}
					messages={messages}
					messageForEdit={messageForEdit}
					selectedMessageIDs={selectedMessageIDs}
					loading={loading}
					isAnyThreadOpen={isAnyOpen}
					checkMyThreadOpen={checkMyThreadOpen}
					onCopyLinkClick={onCopyLinkClick}
					onSelectClick={onSelectClick}
					onUnselectClick={onUnselectClick}
					onEditClick={onEditClick}
					onResetEditClick={reset}
					onToggleThreadClick={onToggleThreadClick}
					onCopyClick={onCopyClick}
				/>
			</Tag>

			<Tag css="w-full shadow-[0_-10px_10px_-10px_rgba(0,0,0,0.4)]">
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
