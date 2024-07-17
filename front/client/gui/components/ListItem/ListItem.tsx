import React from 'react';
import Tag from '../Tag';

const ListItem = props => {
  return (
    <Tag>
      {props.children}
    </Tag>
  )
}

export default ListItem;
