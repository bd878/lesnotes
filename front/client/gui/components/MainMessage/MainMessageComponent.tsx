import React from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is'
import {getFileDownloadUrl} from "../../../api";

function MainMessageComponent(props) {
	const {
		message,
		userID,
	} = props

	if (is.empty(message) || (is.notEmpty(message) && is.empty(message.text)))
		return null

	return (
		<Tag>
			<Tag>
				{message.text}
			</Tag>

			{(message.file && message.file.name && userID) ? <Tag
				el="img"
				src={getFileDownloadUrl(`/files/v2/${userID}/${message.file.name}`)}
			/> : null}
		</Tag>
	)
}

export default MainMessageComponent;