import React, {forwardRef} from 'react';
import Tag from '../Tag';

const List = forwardRef((props, ref) => (
  <Tag el={props.el} ref={ref} css={props.css}>
    {props.children}
  </Tag>
));

export default List;
