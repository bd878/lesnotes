import React from 'react'
import Tag from '../Tag';
import i18n from '../../../i18n';

function MainPageComponent(props) {
	return (
		<Tag css="max-w-md min-w-3xs w-full flex flex-row justify-between">
			<Tag el="a" css="italic text-2xl inline-block cursor-pointer" href="/" target="_self">{i18n("lesnotes")}</Tag>

			<Tag css="inline-block py-1">
				<Tag el="a" href="/login" css="underline">{i18n("login")}</Tag>
				{" | "}
				<Tag el="a" href="/signup" css="underline">{i18n("register")}</Tag>
			</Tag>
		</Tag>
	)
}

export default MainPageComponent;
