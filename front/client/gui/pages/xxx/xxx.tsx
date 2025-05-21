import React from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import Footer from '../../components/Footer';
import i18n from '../../../i18n';

function XXX() {
	return (
		<Tag css="wrap">
			<Tag css="grow">{i18n("not_found")}</Tag>
			<Footer />
		</Tag>
	)
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<XXX />);
