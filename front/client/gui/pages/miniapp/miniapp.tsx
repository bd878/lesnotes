import React from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import * as is from '../../../third_party/is';

function Home() {
	const browser = document.body.dataset.browser
	const isMobile = document.body.dataset.mobile

	return (
		<Tag css="wrap dark">
			{is.notUndef(window.Telegram) ? <Tag>{window.Telegram.WebApp.initData}</Tag> : <Tag>{"window telegram is undef"}</Tag>}
			<Tag css="bg-white dark:bg-black">
				<Footer />
			</Tag>
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Home />);
