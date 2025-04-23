import React from 'react'
import Tag from '../Tag';

function NotificationComponent(props) {
	const {text, css} = props

	return (
		<Tag css={(css || "") + " " + ""}>
			{text}
		</Tag>
	)
}

export default NotificationComponent