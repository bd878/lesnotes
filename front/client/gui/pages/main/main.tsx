import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import StoreProvider from '../../providers/Store';
import MainPage from '../../components/MainPage';
import TgAuth from '../../components/TgAuth';
import AuthProvider from '../../providers/Auth';
import i18n from '../../../i18n';

function Main() {
	return (
		<Tag css="wrap">
			<StoreProvider>
				<AuthProvider inverted={true} shouldFailRedirect={false}>
					<Tag css="w-full grow">
						<MainPage />

						<TgAuth />
					</Tag>
				</AuthProvider>
			</StoreProvider>

			<Footer />
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
