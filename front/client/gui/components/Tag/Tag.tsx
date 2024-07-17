import React from 'react';

const Tag = props => {
  let TagName = props.el || "div";

  return (
    <TagName
      required={props.required}
      href={props.href}
      target={props.target}
      className={props.css}
      name={props.name}
      type={props.type}
      onClick={props.onClick}
      onSubmit={props.onSubmit}
      encType={props.encType}
      action={props.action}
    >{props.children}</TagName>
  )
}

export default Tag;
