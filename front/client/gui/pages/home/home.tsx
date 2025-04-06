import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../../i18n';
import HomePage from '../../components/HomePage';
import AuthProvider from '../../providers/Auth';
import StoreProvider from '../../providers/Store';

function Home() {
	return (
		<StoreProvider>
			<AuthProvider fallback={i18n("messages_auth_fallback")}>
				<HomePage />
			</AuthProvider>
		</StoreProvider>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Home />);
