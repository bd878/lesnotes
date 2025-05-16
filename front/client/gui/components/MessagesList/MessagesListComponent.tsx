import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'
import ArrowRightIcon from '../../icons/ArrowRightIcon'
import ArrowLeftIcon from '../../icons/ArrowLeftIcon'
import LinkIcon from '../../icons/LinkIcon'
import CrossIcon from '../../icons/CrossIcon'
import CrayonIcon from '../../icons/CrayonIcon'
import CopyIcon from '../../icons/CopyIcon'

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"));
const Checkmark = lazy(() => import("../../components/Checkmark"));
const Button = lazy(() => import("../../components/Button"));

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		messages,
		selectedMessageIDs,
		loading,
		error,
		messageForEdit,
		checkMyThreadOpen,
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
		<Tag css={css}>
			{loading ? i18n("loading") : null}
			{error ? null : (
				<List el="ul" css="w-full">
					{messages.map(message => {
						const isMyThreadOpen = is.func(checkMyThreadOpen) ? checkMyThreadOpen(message.ID) : false
						const isSelected = is.notEmpty(selectedMessageIDs) ? selectedMessageIDs.has(message.ID) : false
						const isEdit = is.notEmpty(messageForEdit) ? messageForEdit.ID === message.ID : false
						const isPublic = is.notUndef(message.private) ? !message.private : false

						return (
							<Tag
								el="li"
								tabIndex="0"
								key={`tag_${message.ID}`}
								/* TODO: message.isHovered get computitional property */
								css={
									(liCss || "")
									+ " "
									+ (isAnyThreadOpen ? isMyThreadOpen ? "" : "opacity-50" : "")
									+ " "
									+ "mb-2 ml-8 mr-13 grow-1 max-w-full flex flex-row justify-between"
								}
							>
								<Tag el="label" css="p-2 mt-1 mr-1 cursor-pointer -ml-8">
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
								<Tag css="ml-1 flex -ml-13">
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
				)}
				</List>
			)}
		</Tag>
	)
}

export default MessagesListComponent;
