import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
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
		isPrivate,
		onToggleThreadClick,
		onDeleteClick,
		onPublishClick,
		onPrivateClick,
	} = props

	return (
		<Tag  key={`item_${message.ID}`} el="details" css={(css || "") + " " +  "grow flex flex-col overflow-hidden max-w-full my-1 mr-1 marker:text-xl"}>
			<Tag el="summary" css="cursor-pointer py-1 px-2 hover:bg-gray-300 rounded-sm text-sm italic overflow-hidden whitespace-nowrap text-ellipsis">
				<Tag el="span" css="px-2 py-1">{message.text}</Tag>
			</Tag>

			<Tag css="mt-2">
				{is.trueVal(message.createUTCNano) ? <Tag css="text-sm"><Tag el="span" css="font-bold">{i18n("created_at") + ": "}</Tag>{message.createUTCNano}</Tag> : null}
				{equal(message.updateUTCNano).not(message.createUTCNano) ? <Tag css="text-sm"><Tag el="span" css="font-bold">{i18n("updated_at") + ": "}</Tag>{message.updateUTCNano}</Tag> : null}
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
					{isPrivate
						? is.func(onPublishClick) ? <Button type="button" css="ml-1 btn" onClick={onPublishClick} content={i18n("publish_message")} /> : null
						: is.func(onPrivateClick) ? <Button type="button" css="ml-1 btn" onClick={onPrivateClick} content={i18n("private_message")} /> : null}
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
	)
}

export default MessageElementComponent;
