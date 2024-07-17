import React from 'react';
import Tag from '../Tag';

const FormField = (props) => {
  return (
    <Tag
      el="input"
      name={props.name}
      type={props.type}
      required={props.required}
    />
  );
}

export default FormField;
