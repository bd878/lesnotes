import React, {forwardRef} from 'react';
import Tag from '../Tag';

const FormField = forwardRef((props, ref) => {
	return (
		<Tag
			ref={ref}
			id={props.id}
			css={props.css}
			el={props.el}
			name={props.name}
			value={props.value}
			type={props.type}
			required={props.required}
			onChange={props.onChange}
		/>
	);
});

export default FormField;
