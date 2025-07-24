import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import TgAuth from '../../components/TgAuth';
import StoreProvider from '../../providers/Store';
import AuthProvider from '../../providers/Auth';
import i18n from '../../../i18n';

const LoginForm = lazy(() => import("../../components/LoginForm"));

function Login() {
	return (
		<Tag css="wrap">
			<Suspense fallback={i18n("loading")}>
				<StoreProvider>
					<AuthProvider inverted={true}>
						<Tag css="max-w-md min-w-3xs w-full grow">
							<Tag el="a" css="italic text-2xl mb-8 inline-block cursor-pointer" href="/" target="_self">{i18n("lesnotes")}</Tag>

							<LoginForm />

							<TgAuth />
						</Tag>
					</AuthProvider>
				</StoreProvider>
			</Suspense>

			<Footer />
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
