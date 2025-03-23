import React from 'react';
import Tag from '../Tag';
import * as is from '../../third_party/is'

function MainMessageComponent(props) {
	const {
		message
	} = props

	if (is.empty(message) || (is.notEmpty(message) && is.empty(message.text)))
		return null

	return (
		<Tag>
			{message.text}
		</Tag>
	)
}

export default MainMessageComponent;