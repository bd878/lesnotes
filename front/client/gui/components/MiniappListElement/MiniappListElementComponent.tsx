import React from 'react'
import Tag from '../Tag';
import cn from 'classnames';

function MiniappListElementComponent(props) {
	const {
		message,
		textColor,
		margin,
		radius,
	} = props

	return (
		<Tag el="li" tabIndex="0" css={cn(margin, radius, "bg-body-secondary p-2 overflow-hidden whitespace-nowrap text-ellipsis w-full text-main")}>
			<Tag el="span" css={cn(textColor, "px-2 py-1")}>{message.text}</Tag>
		</Tag>
	)
}

export default MiniappListElementComponent;
