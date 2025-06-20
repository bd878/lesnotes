import React, {useEffect, useState} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import i18n from '../../../i18n';
import api from '../../../api';
import * as is from '../../../third_party/is';

function Miniapp() {
	const [valid, setValid] = useState(null)
	const [loading, setLoading] = useState(true)

	useEffect(() => {
		async function validateData(dataStr) {
			setLoading(true)
			let result = await api.validateMiniappData(dataStr)
			if (result.ok) {
				setValid(true);
			} else {
				setValid(false);
				api.sendLog(JSON.stringify(result))
			}
			setLoading(false)
		}

		validateData(window.Telegram.WebApp.initData)
	}, [setValid, setLoading])

	return (
		<Tag css="wrap dark">
			<Tag css="bg-white">{"Miniapp"}</Tag>
			<Tag css="bg-white">{valid ? "data ok" : "data not ok"}</Tag>
			<Tag css="bg-white dark:bg-black">
				<Footer />
			</Tag>
		</Tag>
	)
}

function Main() {
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
		<Miniapp />
	);
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
