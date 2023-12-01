import React from 'react'

const FormButton = props => {
  return (
    <button type={props.button}>
      {props.text}
    </button>
  )
}

export default FormButton;
