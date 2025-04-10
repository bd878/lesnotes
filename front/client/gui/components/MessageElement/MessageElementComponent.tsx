import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import RubbishIcon from '../../icons/RubbishIcon';
import CrayonIcon from '../../icons/CrayonIcon';
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
			css={css + " " + "cursor-pointer px-2 py-1 bg-gray-100 hover:bg-gray-200 flex flex-row justify-between"}
			key={`tag_${message.ID}`}
			onClick={(e) => {e.stopPropagation(); onClick(message)}} 
		>
			<Tag css="flex flex-row">
				<Tag
					el="button"
					css="cursor-pointer mr-3"
					type="button"
					onClick={(e) => {e.stopPropagation();onEditClick(message)}}
				><CrayonIcon width="18" height="18" /></Tag>

				<Tag css="flex flex-col">
					<ListItem css="mb-1" key={`item_${message.ID}`}>{message.text}</ListItem>

					{(message.file && message.file.ID && message.file.name) ? <Tag
						el="a"
						href={getFileDownloadUrl(`/files/v1/download?id=${message.file.ID}`, false)}
						download={message.file.name}
						target="_blank"
					>
						{message.file.name}
					</Tag> : null}
				</Tag>
			</Tag>

			<Tag
				el="button"
				css="cursor-pointer"
				type="button"
				onClick={(e) => {e.stopPropagation();onDeleteClick(message)}}
			><RubbishIcon css="fill-rose-900/30 hover:fill-rose-900" width="24" height="24" /></Tag>
		</Tag>
	)
}

export default MessageElementComponent;
