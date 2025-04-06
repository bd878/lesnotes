import React, {useEffect} from 'react';
import i18n from '../../../i18n';
import {connect} from '../../../third_party/react-redux';
import {authActionCreator} from '../../features/me'
import {
	selectWillRedirect,
	selectIsAuth,
	selectIsLoading,
} from '../../features/me'

function Auth(props) {
	const {inverted, auth, isAuth, willRedirect, isLoading} = props

	useEffect(() => {auth()}, [auth])

	if (isLoading)
		return (<>{i18n('auth_process')}</>)

	if (willRedirect)
		return (<></>)

	if (inverted)
		return (
			<>{!isAuth
				? props.children
				: (props.fallback || i18n("authed"))
			}</>
		)

	return (
		<>{isAuth
			? props.children
			: (props.fallback || i18n("not_authed"))
		}</>
	);
}

const mapStateToProps = state => ({
	isAuth: selectIsAuth(state),
	isLoading: selectIsLoading(state),
	willRedirect: selectWillRedirect(state),
})

const mapDispatchToProps = ({
	auth: authActionCreator,
})

export default connect(
	mapStateToProps, mapDispatchToProps)(Auth);
