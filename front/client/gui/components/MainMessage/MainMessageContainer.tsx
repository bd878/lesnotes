import React from 'react';
import {connect} from '../../../third_party/react-redux';
import * as is from '../../../third_party/is'
import MainMessageComponent from './MainMessageComponent'
import {
	selectUser,
} from '../../features/me';

function MainMessageContainer(props) {
	const {
		message,
		user,
	} = props

	return (
		<MainMessageComponent message={message} user={user} />
	)
}

const mapStateToProps = state => ({
	user: selectUser(state),
})

const mapDispatchToProps = _dispatch => ({})

export default connect(mapStateToProps, mapDispatchToProps)(MainMessageContainer);