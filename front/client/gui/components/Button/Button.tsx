import React, {forwardRef} from 'react'
import Tag from '../Tag';

const Button = forwardRef((props, ref) => {
	return (
		<Tag
			el="button"
			ref={ref}
			tabIndex={props.tabIndex}
			css={props.css}
			type={props.button}
			onClick={props.onClick}
			disabled={props.disabled}
		>
			{props.content}
		</Tag>
	)
})

export default Button;
