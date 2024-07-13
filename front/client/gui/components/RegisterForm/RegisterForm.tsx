import React, {lazy} from 'react';
import api from '../../api';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form/Form.tsx"));
const FormField = lazy(() => import("../../components/FormField/FormField.tsx"));
const Button = lazy(() => import("../../components/Button/Button.tsx"));

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
}

const RegisterForm = props => {
  return (
    <div>
      <Form name="register-form" onSubmit={sendRegisterRequest}>
        <FormField required name="name" type="text" />
        <FormField required name="password" type="password" />
        <Button type="submit" text={i18n("register")} />
      </Form>
      <a href="/login" target="_self">{i18n("login")}</a>
    </div>
  );
}

export default RegisterForm;
