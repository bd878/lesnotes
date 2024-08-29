import React from 'react';
import Tag from '../Tag';

const ListItem = props => {
  return (
    <Tag el={props.el}>
      {props.children}
    </Tag>
  )
}

export default ListItem;
