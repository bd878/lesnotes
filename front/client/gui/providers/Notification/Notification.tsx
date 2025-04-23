import React, {useEffect} from 'react';
import Tag from '../../components/Tag'
import Notification from '../../components/Notification';
import {connect} from '../../../third_party/react-redux';
import {
	selectIsNotificationVisible,
	selectNotificationText,
} from '../../features/notification'

function NotificationProvider(props) {
	const {children, isVisible, text} = props

	return (
		<>
			{children}
			{isVisible ? (
				<Notification css="absolute" text={text} />
			) : null}
		</>
	)
}

const mapStateToProps = state => ({
	isVisible: selectIsNotificationVisible(state),
	text: selectNotificationText(state),
})

const mapDispatchToProps = ({})

export default connect(mapStateToProps, mapDispatchToProps)(NotificationProvider);