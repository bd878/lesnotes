import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import {getFileDownloadUrl} from "../../../api";

const ListItem = lazy(() => import("../../components/ListItem"));

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
			<Tag css="grow flex flex-row justify-between overflow-hidden max-w-full">
				<Tag css="flex flex-col overflow-hidden mr-3">
					<ListItem css="mb-1 overflow-hidden whitespace-nowrap text-ellipsis" key={`item_${message.ID}`}>{message.text}</ListItem>

					{(message.file && message.file.ID && message.file.name) ? <Tag
						el="a"
						href={getFileDownloadUrl(`/files/v1/download?id=${message.file.ID}`, false)}
						download={message.file.name}
						target="_blank"
					>
						{message.file.name}
					</Tag> : null}
				</Tag>

				{message.createUTCNano ? (
					<Tag css="flex flex-col justify-around whitespace-nowrap">
						<Tag el="div" css="text-xs italic">{i18n("created_at") + ": " + message.createUTCNano}</Tag>
						{message.createUTCNano != message.updateUTCNano ? <Tag el="div" css="text-xs italic">{i18n("updated_at") + ": " + message.updateUTCNano}</Tag> : null}
					</Tag>
				) : null}
			</Tag>
		</Tag>
	)
}

export default MessageElementComponent;
