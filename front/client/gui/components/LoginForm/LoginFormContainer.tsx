import React, {useCallback, useState} from 'react';
import {connect} from '../../third_party/react-redux';
import LoginFormComponent from './LoginFormComponent';
import Tag from '../../components/Tag';
import i18n from '../../i18n';
import {loginActionCreator} from '../../features/me'

function LoginFormContainer(props) {
  const {login} = props

  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const onNameChange = useCallback(e => {
    setName(e.target.value);
  }, [setName]);

  const onPasswordChange = useCallback(e => {
    setPassword(e.target.value);
  }, [setPassword]);

  const sendLoginRequest = useCallback(e => {
    e.preventDefault();

    if (!name) {console.error(i18n("name_required_err")); return;}
    if (!password) {console.error(i18n("pass_required_err")); return;}

    login(name, password)
  }, [login, name, password]);

  return (
    <LoginFormComponent
      name={name}
      password={password}
      onNameChange={onNameChange}
      onPasswordChange={onPasswordChange}
      sendLoginRequest={sendLoginRequest}
    />
  );
}

const mapStateToProps = () => {}

const mapDispatchToProps = ({
  login: loginActionCreator,
})

export default connect(
  mapStateToProps, mapDispatchToProps,
)(LoginFormContainer);
