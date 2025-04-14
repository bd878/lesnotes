import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"))

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		checkMyThreadOpen,
		isAnyThreadOpen,
		onToggleThreadClick,
		onEditClick,
		onDeleteClick,
		onCopyClick,
	} = props

	return (
		<>
			{loading ? i18n("loading") : null}
			{error ? null : (
				<List el="ul" css={css}>
					{messages.map(message => {
						let isMyThreadOpen = is.func(checkMyThreadOpen) ? checkMyThreadOpen(message.ID) : false

						return (
							<Tag
								el="li"
								tabIndex="0"
								/* TODO: message.isHovered get computitional property */
								css={
									(liCss || "")
									+ " "
									+ (isAnyThreadOpen ? isMyThreadOpen ? "" : "opacity-50" : "")
									+ " "
									+ "mb-2 px-2 py-1 mx-1 bg-gray-100 hover:bg-gray-200 flex flex-row justify-between"
								}
								key={`tag_${message.ID}`}
							>
								<MessageElement
									key={message.ID}
									message={message}
									isThreadOpen={isMyThreadOpen}
									onToggleThreadClick={() => onToggleThreadClick(message)}
									onCopyClick={() => onCopyClick(message)}
									onEditClick={() => onEditClick(message)}
									onDeleteClick={() => onDeleteClick(message)}
								/>
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
