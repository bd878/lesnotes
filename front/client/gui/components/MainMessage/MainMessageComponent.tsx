import React from 'react';

function MainMessageComponent(props) {
	const {
		message
	} = props

	return (
		<div>
			{message.text}
		</div>
	)
}

export default MainMessageComponent;