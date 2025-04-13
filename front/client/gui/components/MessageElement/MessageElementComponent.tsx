import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import CrayonIcon from '../../icons/CrayonIcon'
import CopyIcon from '../../icons/CopyIcon'
import * as is from '../../../third_party/is'
import {getFileDownloadUrl} from "../../../api";

function MessageElementComponent(props) {
	const {
		message,
		css,
		tabIndex,
		onOpenThreadClick,
		onEditClick,
		onDeleteClick,
		onCopyClick,
	} = props

	return (
		<Tag
			el="li"
			tabIndex={tabIndex}
			css={css + " " + "cursor-pointer px-2 py-1 mx-1 bg-gray-100 hover:bg-gray-200 flex flex-row justify-between"}
			key={`tag_${message.ID}`}
		>
			<Tag css="grow flex flex-row justify-between items-start overflow-hidden max-w-full">
				<Tag css="flex flex-col overflow-hidden mr-3 w-full">
					<Tag css="mb-1" key={`item_${message.ID}`}>
						<Tag el="details" css="marker:text-xl">
							<Tag el="summary" css="text-sm italic overflow-hidden whitespace-nowrap text-ellipsis">
								<Tag el="span" css="px-2 py-1">{message.text}</Tag>
							</Tag>

							<Tag css="mt-2">
								{message.createUTCNano ? <Tag><Tag el="span" css="font-bold text-sm">{i18n("created_at") + ": "}</Tag>{message.createUTCNano}</Tag> : null}
								{message.updateUTCNano !== message.createUTCNano ? <Tag><Tag el="span" css="font-bold text-sm">{i18n("updated_at") + ": "}</Tag>{message.updateUTCNano}</Tag> : null}
								{is.trueVal(message.fileID) ? (
									<Tag>
										<Tag el="span" css="font-bold text-sm">{i18n("attachments") + ": "}</Tag>
										<Tag el="a" href={getFileDownloadUrl(`/files/v1/download?id={message.fileID}`, false)} target="_blank" download={message.file.name}>{message.file.name}</Tag>
									</Tag>
								) : null}
								<Tag>{message.text}</Tag>
							</Tag>
						</Tag>
					</Tag>
				</Tag>

				<Tag el="button" type="button" onClick={onCopyClick} css="flex pt-[8px]">
					<CopyIcon css="flex ml-1" width="18" height="18" />
				</Tag>
			</Tag>
		</Tag>
	)
}

export default MessageElementComponent;
