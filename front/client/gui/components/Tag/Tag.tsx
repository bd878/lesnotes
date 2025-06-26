import React, {forwardRef} from 'react';

const Tag = forwardRef((props, ref) => {
	let TagName = props.el || "div";

	return (
		<TagName
			ref={ref}
			id={props.id}
			style={props.style}
			tabIndex={props.tabIndex}
			htmlFor={props.htmlFor}
			disabled={props.disabled}
			required={props.required}
			href={props.href}
			target={props.target}
			className={props.css}
			name={props.name}
			type={props.type}
			checked={props.checked}
			value={props.value}
			onClick={props.onClick}
			onChange={props.onChange}
			onScroll={props.onScroll}
			onMouseDown={props.onMouseDown}
			onMouseUp={props.onMouseUp}
			draggable={props.draggable}
			onDragStart={props.onDragStart}
			onDragEnd={props.onDragEnd}
			onDragOver={props.onDragOver}
			onDragEnter={props.onDragEnter}
			onDragLeave={props.onDragLeave}
			onDrop={props.onDrop}
			encType={props.encType}
			action={props.action}
			autoComplete={props.autoComplete}
			download={props.download}
		>{props.children}</TagName>
	)
});

export default Tag;
