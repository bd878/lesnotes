import React from 'react';
import {Provider} from '../../../third_party/react-redux'
import createStore from '../../store';

function StoreProvider(props) {
	return (
		<Provider store={createStore({
			browser: props.browser,
			isMobile: props.isMobile,
			isDesktop: props.isDesktop,
			isMiniapp: props.isMiniapp,
		})}>
			{props.children}
		</Provider>
	)
}

export default StoreProvider;