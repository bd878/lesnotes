import React from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is'
import {getFileDownloadUrl} from "../../../api";

function MainMessageComponent(props) {
	const {
		message
	} = props

	if (is.empty(message) || (is.notEmpty(message) && is.empty(message.text)))
		return null

	return (
		<Tag>
			<Tag>
				{message.text}
			</Tag>

			{(message.file && message.file.ID && message.file.name) ? <Tag
				el="a"
				href={getFileDownloadUrl(`/files/v1/download?id=${message.file.ID}`, false)}
				download={message.file.name}
				target="_blank"
			>
				{message.file.name}
			</Tag> : null}
		</Tag>
	)
}

export default MainMessageComponent;