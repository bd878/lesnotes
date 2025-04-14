import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is'

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"))

function MessagesListComponent(props) {
	const {
		css,
		messages,
		loading,
		error,
		threadID,
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
					{messages.map(message => (
						<Tag
							el="li"
							tabIndex="0"
							/* TODO: message.isHovered get computitional property */
							css={((is.notEmpty(threadID) && threadID !== message.ID) ? "opacity-50 " : "") + "mb-2 px-2 py-1 mx-1 bg-gray-100 hover:bg-gray-200 flex flex-row justify-between"}
							key={`tag_${message.ID}`}
						>
							<MessageElement
								key={message.ID}
								message={message}
								isThreadOpen={is.notEmpty(threadID) && threadID === message.ID}
								onToggleThreadClick={() => onToggleThreadClick(message)}
								onCopyClick={() => onCopyClick(message)}
								onEditClick={() => onEditClick(message)}
								onDeleteClick={() => onDeleteClick(message)}
							/>
						</Tag>
					))}
				</List>
			)}
		</>
	)
}

export default MessagesListComponent;
