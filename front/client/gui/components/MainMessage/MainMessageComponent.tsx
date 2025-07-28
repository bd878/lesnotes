import React from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is'
import {getFileDownloadUrl} from "../../../api";

function MainMessageComponent(props) {
	const {
		message,
		user,
	} = props

	if (is.empty(message) || (is.notEmpty(message) && is.empty(message.text)))
		return null

	return (
		<Tag>
			<Tag>
				{message.text}
			</Tag>

			{(message.file && message.file.name && user && user.ID) ? <Tag
				el="img"
				src={getFileDownloadUrl(`/files/v2/${user.ID}/${message.file.name}`)}
			/> : null}
		</Tag>
	)
}

export default MainMessageComponent;