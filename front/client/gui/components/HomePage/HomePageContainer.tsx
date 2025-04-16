import React, {useCallback} from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../../third_party/react-redux';
import {
	selectStack,
	openThreadActionCreator,
	closeThreadActionCreator,
} from '../../features/stack';

function HomePageContainer(props) {
	const {
		stack,
		open,
		close,
	} = props

	const openThread = useCallback(index => message => open({index, threadID: message.ID}), [open])
	const closeThread = useCallback(index => _message => close({index}), [close])

	return (
		<HomePageComponent
			stack={stack}
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
}

export default connect(mapStateToProps, mapDispatchToProps)(HomePageContainer);
