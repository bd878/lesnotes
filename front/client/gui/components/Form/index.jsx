import React from 'react';

const Form = props => {
  return (
    <form action="" name={props.name} onSubmit={props.onSubmit} encType={props.enctype||null}>
      {props.children}
    </form>
  )
}

export default Form;
