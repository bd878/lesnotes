import React, {forwardRef} from 'react'
import Tag from '../Tag';

const Checkmark = forwardRef((props, ref) => {
	return (
		<Tag
			el="input"
			type="checkbox"
		></Tag>
	)
})

export default Checkmark;