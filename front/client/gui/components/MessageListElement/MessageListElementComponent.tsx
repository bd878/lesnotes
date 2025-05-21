import React, {lazy} from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is'
import i18n from '../../../i18n';
import ArrowRightIcon from '../../icons/ArrowRightIcon'
import ArrowLeftIcon from '../../icons/ArrowLeftIcon'
import LinkIcon from '../../icons/LinkIcon'
import CrossIcon from '../../icons/CrossIcon'
import CrayonIcon from '../../icons/CrayonIcon'
import CopyIcon from '../../icons/CopyIcon'
import DotsIcon from '../../icons/DotsIcon'

const Checkmark = lazy(() => import("../../components/Checkmark"));
const Button = lazy(() => import("../../components/Button"));
const MessageElement = lazy(() => import("../../components/MessageElement"));

function MessageListElementComponent(props) {
	const {
		css,
		message,
		listRef,
		onDragStart,
		isMyThreadOpen,
		isSelected,
		isEdit,
		isPublic,
		isAnyThreadOpen,
		onToggleThreadClick,
		onEditClick,
		onCopyClick,
		onCopyLinkClick,
		onSelectClick,
		onUnselectClick,
		onResetEditClick,
	} = props

	return (
		<Tag
			el="li"
			tabIndex="0"
			ref={listRef}
			/* TODO: message.isHovered get computitional property */
			css={
				(css || "")
				+ " "
				+ (isAnyThreadOpen ? isMyThreadOpen ? "" : "opacity-50" : "")
				+ " "
				+ "pb-2 pr-[80px] grow-1 max-w-full flex flex-row justify-between cursor-move"
			}
		>
			<Tag el="label" css="p-2 mt-1 mr-1 cursor-pointer -pl-8">
				<Checkmark
					css="cursor-pointer"
					tabIndex="0"
					onChange={() => isSelected ? onUnselectClick(message) : onSelectClick(message)}
					name=""
					id={message.ID}
					value={message.ID}
					checked={isSelected}
				/>
			</Tag>
			<Tag css="px-2 py-1 bg-gray-100 grow-1 flex flex-row items-start justify-between max-w-full min-w-full">
				<Tag draggable onDragStart={onDragStart} css="mr-1 my-[10px] cursor-pointer flex items-center">
					<DotsIcon width="24" height="24" />
				</Tag>

				<MessageElement
					key={message.ID}
					message={message}
					isPrivate={message.private}
					isThreadOpen={isMyThreadOpen}
				/>

				{isPublic
					? is.func(onCopyLinkClick)
						? (
							<Button
								type="button"
								css="flex my-1 p-2 rounded-sm cursor-pointer hover:bg-gray-300"
								content={
									<LinkIcon css="flex" width="20" height="20" />
								}
								onClick={() => onCopyLinkClick(message)}
							/>
						) : null
					: null
				}
				{isEdit
					? is.func(onResetEditClick)
						? (
							<Button
								type="button"
								css="flex my-1 p-2 rounded-sm cursor-pointer hover:bg-gray-300"
								content={
									<CrossIcon css="flex" width="20" height="20" />
								}
								onClick={() => onResetEditClick()}
							/>
						) : null

					: is.func(onEditClick)
						? (
							<Button
								type="button"
								css="flex my-1 p-2 rounded-sm cursor-pointer hover:bg-gray-300"
								content={
									<CrayonIcon css="flex" width="20" height="20" />
								}
								onClick={() => onEditClick(message)}
							/>
						) : null
				}
				{is.func(onCopyClick) ? (
					<Button
						type="button"
						css="flex my-1 p-2 rounded-sm cursor-pointer hover:bg-gray-300"
						content={
							<CopyIcon css="flex" width="20" height="20" />
						}
						onClick={() => onCopyClick(message)}
					/>
				) : null}
			</Tag>
			<Tag css="ml-1 flex -pl-13">
				<Button
					type="button"
					tabIndex="0"
					css={(isMyThreadOpen ? "bg-gray-100 " : "") + "h-[52px] hover:bg-gray-300 p-2 cursor-pointer rounded-sm flex items-center"}
					content={
						isMyThreadOpen ? <ArrowLeftIcon width="24" height="24" /> : <ArrowRightIcon width="24" height="24" />
					}
					onClick={() => onToggleThreadClick(message)}
				/>
			</Tag>
		</Tag>
	)
}

export default MessageListElementComponent;
