import React, {Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import MessagePage from '../../components/MessagePage';
import Footer from '../../components/Footer';
import StoreProvider from '../../providers/Store';

function Message() {
	const id = document.body.dataset.id

	return (
		<Tag css="wrap">
			<StoreProvider>
				<Tag css="grow w-full">
					<MessagePage id={id} />
				</Tag>
			</StoreProvider>

			<Footer />
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Message />);
