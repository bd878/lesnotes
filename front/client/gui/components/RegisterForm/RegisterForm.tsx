import React, {lazy} from 'react';
import api from '../../api';
import Tag from '../Tag';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

const sendRegisterRequest = async e => {
  e.preventDefault();

  let name = e.target.elements['name'].value;
  let password = e.target.elements['password'].value;
  if (!name) {console.error(i18n("name_required_err")); return;}
  if (!password) {console.error(i18n("pass_required_err")); return;}

  const response = await api("/users/v1/signup", {
    method: "POST",
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
    },
    body: new URLSearchParams({
      'name': name,
      'password': password,
    })
  });
  console.log(response);
  if (response.error != "") {
    console.error("[RegisterForm]: /signup response returned error", response.error, response.explain)
  } else {
    if (response.value.status == "ok") {
      setTimeout(() => {location.href = "/home"}, 1000)
    }
  }
}

const RegisterForm = props => {
  return (
    <Tag>
      <Form name="register-form" onSubmit={sendRegisterRequest}>
        <FormField required name="name" type="text" />
        <FormField required name="password" type="password" />
        <Button type="submit" text={i18n("register")} />
      </Form>
      <Tag el="a" href="/login" target="_self">{i18n("login")}</Tag>
    </Tag>
  );
}

export default RegisterForm;
