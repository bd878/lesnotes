import React, {useEffect} from 'react';
import i18n from '../../../i18n';
import {connect} from '../../../third_party/react-redux';
import {authActionCreator} from '../../features/me'
import {
	selectWillRedirect,
	selectIsAuth,
	selectIsLoading,
} from '../../features/me'

function AuthProvider(props) {
	const {inverted, auth, isAuth, willRedirect, isLoading} = props

	useEffect(() => {auth()}, [auth])

	if (isLoading)
		return i18n("loading")

	if (willRedirect)
		return (<></>)

	if (inverted)
		return (
			<>{!isAuth
				? props.children
				: (props.fallback || <Tag css="m-8 mt-10">{i18n("authed")}</Tag>)
			}</>
		)

	return (
		<>{isAuth
			? props.children
			: (props.fallback || <Tag css="m-8 mt-10">{i18n("not_authed")}</Tag>)
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
	mapStateToProps, mapDispatchToProps)(AuthProvider);
