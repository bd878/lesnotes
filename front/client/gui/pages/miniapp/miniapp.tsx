import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom/client';
import {init, backButton} from '@telegram-apps/sdk-react';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import i18n from '../../../i18n';
import api from '../../../api';
import * as is from '../../../third_party/is';

init();

backButton.mount();

function Home() {
	const [valid, setValid] = useState(null)
	const [loading, setLoading] = useState(true)

	const browser = document.body.dataset.browser
	const isMobile = document.body.dataset.mobile

	useEffect(() => {
		if (is.undef(window.Telegram))
			return

		async function validateData(dataStr) {
			setLoading(true)
			let result = await api.validateMiniappData(dataStr)
			if (result.ok)
				setValid(true);
			else
				setValid(false);
			setLoading(false)
		}

		validateData(window.Telegram.WebApp.initData)
	}, [window.Telegram, setValid, setLoading])

	if (is.undef(window.Telegram))
		return (
			<Tag css="wrap">
				{i18n("miniapp_only")}
				<Tag css="bg-white">
					<Footer />
				</Tag>
			</Tag>
		);

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
