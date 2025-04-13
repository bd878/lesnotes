import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import StoreProvider from '../../providers/Store';

const RegisterForm = lazy(() => import("../../components/RegisterForm"));

function Register() {
	return (
		<Tag css="wrap">
			<Suspense fallback={i18n("loading")}>
				<StoreProvider>
					<Tag css="max-w-md min-w-3xs w-full">
						<Tag el="a" css="italic text-2xl mb-8 inline-block cursor-pointer" href="/" target="_self">{i18n("lesnotes")}</Tag>

						<RegisterForm />
					</Tag>
				</StoreProvider>
			</Suspense>
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Register />);
