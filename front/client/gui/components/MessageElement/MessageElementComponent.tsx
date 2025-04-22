import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import CrayonIcon from '../../icons/CrayonIcon'
import CopyIcon from '../../icons/CopyIcon'
import * as is from '../../../third_party/is'
import {getFileDownloadUrl} from "../../../api";
import {equal} from '../../../utils';

const Button = lazy(() => import("../../components/Button"));

function MessageElementComponent(props) {
	const {
		message,
		css,
		tabIndex,
		isThreadOpen,
		onToggleThreadClick,
		onEditClick,
		onDeleteClick,
		onCopyClick,
	} = props

	return (
		<Tag css={(css || "") + " " + "grow flex flex-row justify-between items-start overflow-hidden max-w-full"}>
			<Tag css="flex flex-col overflow-hidden w-full">
				<Tag key={`item_${message.ID}`}>
					<Tag el="details" css="m-1 marker:text-xl">
						<Tag el="summary" css="cursor-pointer py-1 px-2 hover:bg-gray-300 rounded-sm text-sm italic overflow-hidden whitespace-nowrap text-ellipsis">
							<Tag el="span" css="px-2 py-1">{message.text}</Tag>
						</Tag>

						<Tag css="mt-2">
							{is.trueVal(message.createUTCNano) ? <Tag><Tag el="span" css="font-bold text-sm">{i18n("created_at") + ": "}</Tag>{message.createUTCNano}</Tag> : null}
							{equal(message.updateUTCNano).not(message.createUTCNano) ? <Tag><Tag el="span" css="font-bold text-sm">{i18n("updated_at") + ": "}</Tag>{message.updateUTCNano}</Tag> : null}
							{is.trueVal(message.fileID) ? (
								<Tag>
									<Tag el="span" css="font-bold text-sm">{i18n("attachments") + ": "}</Tag>
									<Tag
										el="a"
										css="underline text-blue-600 visited:text-purple-600"
										href={getFileDownloadUrl(`/files/v1/download?id=${message.fileID}`, false)}
										target="_blank"
										download={message.file.name}
									>
										{message.file.name}
									</Tag>
								</Tag>
							) : null}
							<Tag>{message.text}</Tag>

							<Tag css="flex flex-row items-start mt-2">
								{is.func(onDeleteClick) ? <Button type="button" css="btn" onClick={onDeleteClick} content={i18n("delete_message")} /> : null}
								{is.func(onEditClick) ? <Button type="button" css="ml-1 btn" onClick={onEditClick} content={i18n("edit_message")} /> : null}
								{is.func(onToggleThreadClick) ? (
									<Button
										type="button"
										css="ml-1 btn"
										onClick={onToggleThreadClick}
										content={
											isThreadOpen
												? i18n("close_thread") + " X"
												: i18n("open_thread") + " >"
										}
									/>
								) : null}
							</Tag>
						</Tag>
					</Tag>
				</Tag>
			</Tag>

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

export default MessageElementComponent;
