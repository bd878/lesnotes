import React from 'react';
import Tag from '../Tag';

function ListItem(props){
  const {
    el,
    onClick,
  } = props

  return (
    <Tag el={el} onClick={onClick}>
      {props.children}
    </Tag>
  )
}

export default ListItem;
