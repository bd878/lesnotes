import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../../i18n';
import Tag from '../../components/Tag';
import HomePage from '../../components/HomePage';
import AuthProvider from '../../providers/Auth';
import StoreProvider from '../../providers/Store';
import * as is from '../../../third_party/is';

function Home() {
	const browser = document.body.dataset.browser
	const isMobile = document.body.dataset.mobile

	return (
		<Tag css="wrap">
			<StoreProvider browser={browser} isMobile={is.trueVal(isMobile)} isDesktop={true}>
				<AuthProvider fallback={i18n("messages_auth_fallback")}>
					<HomePage />
				</AuthProvider>
			</StoreProvider>
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Home />);
