import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'

const List = lazy(() => import("../../components/List"));
const MessageListElement = lazy(() => import("../../components/MessageListElement"));

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		index,
		messages,
		selectedMessageIDs,
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
			{error ? null : (
				<List el="ul" css="w-full">
					{messages.map(message => {
						const isMyThreadOpen = is.func(checkMyThreadOpen) ? checkMyThreadOpen(message.ID) : false
						const isSelected = is.notEmpty(selectedMessageIDs) ? selectedMessageIDs.has(message.ID) : false
						const isEdit = is.notEmpty(messageForEdit) ? messageForEdit.ID === message.ID : false
						const isPublic = is.notUndef(message.private) ? !message.private : false

						return (
							<MessageListElement
								key={`tag_${message.ID}`}
								css={liCss}
								index={index}
								message={message}
								isMyThreadOpen={isMyThreadOpen}
								isSelected={isSelected}
								isEdit={isEdit}
								isPublic={isPublic}
								isAnyThreadOpen={isAnyThreadOpen}
								onToggleThreadClick={onToggleThreadClick}
								onEditClick={onEditClick}
								onCopyClick={onCopyClick}
								onCopyLinkClick={onCopyLinkClick}
								onSelectClick={onSelectClick}
								onUnselectClick={onUnselectClick}
								onResetEditClick={onResetEditClick}
							/>
						)
					}
				)}
				</List>
			)}
		</Tag>
	)
}

export default MessagesListComponent;
