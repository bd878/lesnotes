import React from 'react';

const FormField = (props) => {
  return (
    <input
      name={props.name}
      type={props.type}
    />
  );
}

export default FormField;
