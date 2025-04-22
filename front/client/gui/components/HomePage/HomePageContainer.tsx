import React, {useCallback} from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../../third_party/react-redux';
import {
	selectStack,
	destroyThreadActionCreator,
	openThreadActionCreator,
	closeThreadActionCreator,
} from '../../features/stack';
import {logoutActionCreator} from '../../features/me';

function HomePageContainer(props) {
	const {
		stack,
		open,
		close,
		destroy,
		logout,
	} = props

	const openThread = useCallback(index => message => open({index, threadID: message.ID}), [open])
	const closeThread = useCallback(index => _message => close({index}), [close])
	const destroyThread = useCallback(index => () => index === 0 ? logout() : destroy({index}), [destroy, logout])

	return (
		<HomePageComponent
			stack={stack}
			destroyThread={destroyThread}
			openThread={openThread}
			closeThread={closeThread}
		/>
	)
}

const mapStateToProps = state => ({
	stack: selectStack(state),
})

const mapDispatchToProps = {
	logout: logoutActionCreator,
	open: openThreadActionCreator,
	close: closeThreadActionCreator,
	destroy: destroyThreadActionCreator,
}

export default connect(mapStateToProps, mapDispatchToProps)(HomePageContainer);
