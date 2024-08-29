import React from 'react';
import Tag from '../Tag';

const List = props => (
  <Tag el={props.el} css={props.css}>
    {props.children}
  </Tag>
);

export default List;
