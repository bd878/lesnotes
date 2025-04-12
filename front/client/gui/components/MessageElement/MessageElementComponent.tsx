import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import CrayonIcon from '../../icons/CrayonIcon'
import CopyIcon from '../../icons/CopyIcon'
import {getFileDownloadUrl} from "../../../api";

function MessageElementComponent(props) {
	const {
		message,
		css,
		onClick,
		tabIndex,
		onEditClick,
		onDeleteClick,
	} = props

	return (
		<Tag
			el="li"
			tabIndex={tabIndex}
			css={css + " " + "cursor-pointer px-2 py-1 mx-1 bg-gray-100 hover:bg-gray-200 flex flex-row justify-between"}
			key={`tag_${message.ID}`}
			onClick={(e) => {e.stopPropagation(); onClick(message)}} 
		>
			<Tag css="grow flex flex-row justify-between items-start overflow-hidden max-w-full">
				<Tag css="flex flex-col overflow-hidden mr-3">
					<Tag css="mb-1" key={`item_${message.ID}`}>
						<Tag el="details">
							<Tag el="summary" css="text-sm italic overflow-hidden whitespace-nowrap text-ellipsis">
								{message.text}
							</Tag>

							<Tag css="mt-2">{message.text}</Tag>
						</Tag>
					</Tag>
				</Tag>

				<Tag el="span" css="flex pt-[3px]">
					<CopyIcon css="flex ml-1" width="18" height="18" />
				</Tag>
			</Tag>
		</Tag>
	)
}

export default MessageElementComponent;
