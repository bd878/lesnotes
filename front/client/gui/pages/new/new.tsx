import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';
import NewThreadPage from '../../components/NewThreadPage';
import StoreProvider from '../../providers/Store';

function Message() {
	return (
		<StoreProvider>
			<NewThreadPage />
		</StoreProvider>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Message />);
