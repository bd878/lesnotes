import React from 'react';

const Form = props => {
  return (
    <form action="" name={props.name} onSubmit={props.onSubmit}>
      {props.children}
    </form>
  )
}

export default Form;
