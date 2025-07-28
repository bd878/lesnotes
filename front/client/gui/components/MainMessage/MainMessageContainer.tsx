import React from 'react';
import {connect} from '../../../third_party/react-redux';
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
		<MainMessageComponent message={message} userID={user.ID} />
	)
}

const mapStateToProps = state => ({
	user: selectUser(state),
})

const mapDispatchToProps = _dispatch => ({})

export default connect(mapStateToProps, mapDispatchToProps)(MainMessageContainer);