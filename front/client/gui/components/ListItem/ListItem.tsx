import React from 'react';
import Tag from '../Tag';

function ListItem(props){
	const {
		el,
		onClick,
	} = props

	return (
		<Tag css={props.css} el={el} onClick={onClick}>
			{props.children}
		</Tag>
	)
}

export default ListItem;
