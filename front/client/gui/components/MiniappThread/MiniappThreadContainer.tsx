import React, {useEffect, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';

function MiniappThreadContainer(props) {
	const {
		thread,
	} = props

	return (
		<Tag>
			{thread.list.map((message, index) => (
				<Tag css="px-2 py-1">{message.text}</Tag>
			))}
		</Tag>
	)
}

export default MiniappThreadContainer;
