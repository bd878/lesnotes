import React, {forwardRef} from 'react';

const Tag = forwardRef((props, ref) => {
	let TagName = props.el || "div";

	return (
		<TagName
			ref={ref}
			id={props.id}
			htmlFor={props.htmlFor}
			disabled={props.disabled}
			required={props.required}
			href={props.href}
			target={props.target}
			className={props.css}
			name={props.name}
			type={props.type}
			value={props.value}
			onClick={props.onClick}
			onChange={props.onChange}
			onScroll={props.onScroll}
			encType={props.encType}
			action={props.action}
			autoComplete={props.autoComplete}
			download={props.download}
		>{props.children}</TagName>
	)
});

export default Tag;
