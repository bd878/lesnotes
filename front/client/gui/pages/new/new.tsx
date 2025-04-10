import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import NewThreadPage from '../../components/NewThreadPage';
import StoreProvider from '../../providers/Store';

function Message() {
	return (
		<Tag css="wrap">
			<StoreProvider>
				<NewThreadPage />
			</StoreProvider>
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Message />);
