import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import MessagePage from '../../components/MessagePage';
import StoreProvider from '../../providers/Store';

function Message() {
	const id = document.body.dataset.id

	return (
		<Tag css="wrap">
			<StoreProvider>
				<MessagePage id={id} />
			</StoreProvider>
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Message />);
