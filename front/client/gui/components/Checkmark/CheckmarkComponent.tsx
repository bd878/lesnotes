import React, {forwardRef} from 'react'
import Tag from '../Tag';

const CheckmarkComponent = forwardRef((props, ref) => {
	return (
		<Tag
			css={props.css}
			tabIndex={props.tabIndex}
			el="input"
			type="checkbox"
			value={props.value}
			id={props.id}
			name={props.name}
		></Tag>
	)
})

export default CheckmarkComponent;