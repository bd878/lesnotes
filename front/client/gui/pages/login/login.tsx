import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import StoreProvider from '../../providers/Store';
import AuthProvider from '../../providers/Auth';
import i18n from '../../../i18n';

const LoginForm = lazy(() => import("../../components/LoginForm"));

function Login() {
	return (
		<Suspense fallback={i18n('loading')}>
			<StoreProvider>
				<AuthProvider inverted={true}>
					<Tag css="m-8 mt-10 max-w-md min-w-3xs">
						<Tag css="italic text-2xl mb-8">{i18n("lesnotes")}</Tag>

						<LoginForm />
					</Tag>
				</AuthProvider>
			</StoreProvider>
		</Suspense>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
