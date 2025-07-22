import React from 'react';
import Tag from '../Tag';

function FooterComponent(props) {
	return (
		<Tag css="w-full flex justify-center">
			<Tag css="text-sm text-gray-950 dark:text-white text-(--tg-theme-text-color)">Les notes © 2025</Tag>
		</Tag>
	)
}

export default FooterComponent