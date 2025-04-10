import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';

function Main() {
	return (
		<Tag css="wrap">
			<Tag css="max-w-md min-w-3xs w-full flex flex-row justify-between">
				<Tag el="a" css="italic text-2xl inline-block cursor-pointer" href="/" target="_self">{i18n("lesnotes")}</Tag>

				<Tag css="inline-block py-1">
					<Tag el="a" href="/login" css="underline">{i18n("login")}</Tag>
					{" | "}
					<Tag el="a" href="/signup" css="underline">{i18n("register")}</Tag>
				</Tag>
			</Tag>

		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
