import React from 'react';
import HomePageComponent from './HomePageComponent';
import {connect} from '../../../third_party/react-redux';
import {
	selectStack,
} from '../../features/stack';

function HomePageContainer(props) {
	const {
		stack,
	} = props

	return (
		<HomePageComponent stack={stack} />
	)
}

const mapStateToProps = state => ({
	stack: selectStack(state),
})

const mapDispatchToProps = {
}

export default connect(mapStateToProps, mapDispatchToProps)(HomePageContainer);
