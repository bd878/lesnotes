import React from 'react'
import Tag from '../Tag';

const Button = props => {
  return (
    <Tag el="button" type={props.button} onClick={props.onClick}>
      {props.text}
    </Tag>
  )
}

export default Button;
