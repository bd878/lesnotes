import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"));
const Checkmark = lazy(() => import("../../components/Checkmark"));
const Button = lazy(() => import("../../components/Button"));
const CopyIcon = lazy(() => import('../../icons/CopyIcon'));
const CrayonIcon = lazy(() => import('../../icons/CrayonIcon'));
const CrossIcon = lazy(() => import('../../icons/CrossIcon'));

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
		onPublishClick,
		onPrivateClick,
		onSelectClick,
		onUnselectClick,
		onResetEditClick,
	} = props

	return (
		<>
			{loading ? i18n("loading") : null}
			{error ? null : (
				<List el="ul" css={css}>
					{messages.map(message => {
						const isMyThreadOpen = is.func(checkMyThreadOpen) ? checkMyThreadOpen(message.ID) : false
						const isSelected = is.notEmpty(selectedMessageIDs) ? selectedMessageIDs.has(message.ID) : false
						const isEdit = is.notEmpty(messageForEdit) ? messageForEdit.ID === message.ID : false

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
									+ "mb-2 px-2 py-1 bg-gray-100 grow-1 max-w-full flex flex-row items-start justify-between ml-8"
								}
							>
								<Tag css="mr-5 -ml-8 mt-2">
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
								<MessageElement
									key={message.ID}
									message={message}
									isThreadOpen={isMyThreadOpen}
									isPrivate={message.private}
									onToggleThreadClick={() => onToggleThreadClick(message)}
									onEditClick={() => onEditClick(message)}
									onPublishClick={() => onPublishClick(message)}
									onPrivateClick={() => onPrivateClick(message)}
								/>

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
						)
					}
				)}
				</List>
			)}
		</>
	)
}

export default MessagesListComponent;
