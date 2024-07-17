import React from 'react';
import Tag from '../Tag';

const Form = props => (
  <Tag
    el="form"
    css={props.css}
    action=""
    name={props.name}
    onSubmit={props.onSubmit}
    encType={props.enctype||null}
  >
    {props.children}
  </Tag>
);

export default Form;
