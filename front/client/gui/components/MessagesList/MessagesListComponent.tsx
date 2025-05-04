import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"));
const Checkmark = lazy(() => import("../../components/Checkmark"));
const Button = lazy(() => import("../../components/Button"));
const CopyIcon = lazy(() => import('../../icons/CopyIcon'))

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		messages,
		selectedMessageIDs,
		loading,
		error,
		checkMyThreadOpen,
		isAnyThreadOpen,
		onToggleThreadClick,
		onEditClick,
		onDeleteClick,
		onCopyClick,
		onSelectClick,
		onUnselectClick,
	} = props

	return (
		<>
			{loading ? i18n("loading") : null}
			{error ? null : (
				<List el="ul" css={css}>
					{messages.map(message => {
						let isMyThreadOpen = is.func(checkMyThreadOpen) ? checkMyThreadOpen(message.ID) : false

						const selected = selectedMessageIDs.has(message.ID)

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
								<Checkmark
									css="cursor-pointer mr-5 -ml-8"
									tabIndex="0"
									onClick={() => selected ? onUnselectClick(message) : onSelectClick(message)}
									name=""
									id={message.ID}
									value={message.ID}
								/>
								<MessageElement
									key={message.ID}
									message={message}
									isThreadOpen={isMyThreadOpen}
									onToggleThreadClick={() => onToggleThreadClick(message)}
									onEditClick={() => onEditClick(message)}
									onDeleteClick={() => onDeleteClick(message)}
								/>

								{is.func(onCopyClick) ? (
									<Button
										type="button"
										css="flex my-1 p-2 rounded-sm cursor-pointer hover:bg-gray-300"
										content={
											<CopyIcon css="flex" width="20" height="20" />
										}
										onClick={onCopyClick}
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
