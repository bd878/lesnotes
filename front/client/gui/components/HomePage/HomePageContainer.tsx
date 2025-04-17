import React, {useCallback} from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../../third_party/react-redux';
import {
	selectStack,
	destroyThreadActionCreator,
	openThreadActionCreator,
	closeThreadActionCreator,
} from '../../features/stack';

function HomePageContainer(props) {
	const {
		stack,
		open,
		close,
		destroy,
	} = props

	const openThread = useCallback(index => message => open({index, threadID: message.ID}), [open])
	const closeThread = useCallback(index => _message => close({index}), [close])
	const destroyThread = useCallback(index => () => destroy({index}), [destroy])

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
	open: openThreadActionCreator,
	close: closeThreadActionCreator,
	destroy: destroyThreadActionCreator,
}

export default connect(mapStateToProps, mapDispatchToProps)(HomePageContainer);
