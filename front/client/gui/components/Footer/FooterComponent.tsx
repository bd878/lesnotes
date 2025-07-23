import React from 'react';
import Tag from '../Tag';
import cn from 'classnames';
import * as is from '../../../third_party/is'

function FooterComponent(props) {
	const {
		textColor,
	} = props

	const css = cn(textColor, "text-sm", {
		"text-gray-950 dark:text-white": !textColor,
	})

	return (
		<Tag css="w-full flex justify-center">
			<Tag css={css}>Les Notes © 2025</Tag>
		</Tag>
	)
}

export default FooterComponent