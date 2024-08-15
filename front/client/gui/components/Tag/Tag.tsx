import React, {forwardRef} from 'react';

const Tag = forwardRef((props, ref) => {
  let TagName = props.el || "div";

  return (
    <TagName
      ref={ref}
      required={props.required}
      href={props.href}
      target={props.target}
      className={props.css}
      name={props.name}
      type={props.type}
      value={props.value}
      onClick={props.onClick}
      onChange={props.onChange}
      encType={props.encType}
      action={props.action}
      autoComplete={props.autoComplete}
    >{props.children}</TagName>
  )
});

export default Tag;
