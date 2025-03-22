import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';
import MessagePage from '../../components/MessagePage';
import StoreProvider from '../../providers/Store';

function Message() {
	const id = document.body.dataset.Id
	return (
		<StoreProvider>
			<MessagePage id={id} />
		</StoreProvider>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Message />);
