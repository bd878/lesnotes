import React, {useEffect} from 'react';
import i18n from '../../../i18n';
import Tag from '../../components/Tag'
import {connect} from '../../../third_party/react-redux';
import {authActionCreator} from '../../features/me'
import {
	selectWillRedirect,
	selectIsAuth,
	selectIsLoading,
} from '../../features/me'

function AuthProvider(props) {
	const {inverted, auth, isAuth, willRedirect, isLoading, fallback, children} = props

	useEffect(() => {auth()}, [auth])

	let shouldAllow = inverted ? !isAuth : isAuth

	return (
		<>
			{isLoading
				? i18n("loading")
				: willRedirect
					? null
					: shouldAllow
						? children
						: fallback || <Tag css="m-8 mt-10">{i18n("not_authed")}</Tag>
			}
		</>
	)
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
