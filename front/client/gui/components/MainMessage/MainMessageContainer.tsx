import React from 'react';
import MainMessageComponent from './MainMessageComponent'

function MainMessageContainer(props) {
	const {
		message
	} = props

	return (
		<MainMessageComponent message={message} />
	)
}

export default MainMessageContainer;