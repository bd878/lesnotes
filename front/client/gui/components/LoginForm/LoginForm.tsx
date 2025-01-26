import React, {lazy, useCallback, useState} from 'react';
import api from '../../api';
import Tag from '../../components/Tag';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

const LoginForm = ({ onError }) => {
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

    const send = async (name, password) => {
      try {
        const result = await api.login(name, password)
        if (result.isOk) {
          setTimeout(() => {location.href = "/home"}, 0)
        } else {
          console.log("login result not ok")
          console.error(result.explain)
        }
      } catch (e) {
        console.error(e)
      }
    }

    if (!name) {console.error(i18n("name_required_err")); return;}
    if (!password) {console.error(i18n("pass_required_err")); return;}

    send(name, password)
  }, [
    name,
    password,
    setName,
    setPassword,
  ]);

  return (
    <>
      <Form
        autoComplete="off"
        name="login-form"
      >
        <FormField
          required
          el="input"
          name="name"
          type="text"
          value={name}
          onChange={onNameChange}
        />
        <FormField
          required
          el="input"
          name="password"
          type="password"
          value={password}
          onChange={onPasswordChange}
        />
      </Form>

      <Button
        type="button"
        text={i18n("login")}
        onClick={sendLoginRequest}
      />

      <Tag
        el="a"
        href="/register"
        target="_self"
      >
        {i18n("register")}
      </Tag>
    </>
  );
}

export default LoginForm;
