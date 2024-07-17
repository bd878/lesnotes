import React from 'react';
import Tag from '../Tag';

const List = props => (
  <Tag css={props.css}>
    {props.children}
  </Tag>
);

export default List;
