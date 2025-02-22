import React, {useEffect} from 'react';
import i18n from '../../i18n';
import {connect} from '../../third_party/react-redux';
import {authActionCreator} from '../../features/me'
import {selectIsAuth, selectIsLoading} from '../../features/me'

function Auth(props) {
  const {auth, isAuth, isLoading} = props

  useEffect(() => {auth()}, [auth])

  if (isLoading) {
    return (<>{i18n('auth_process')}</>)
  }

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
})

const mapDispatchToProps = ({
  auth: authActionCreator,
})

export default connect(
  mapStateToProps, mapDispatchToProps)(Auth);
