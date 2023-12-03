import React from 'react'

const Button = props => {
  return (
    <button type={props.button}>
      {props.text}
    </button>
  )
}

export default Button;
