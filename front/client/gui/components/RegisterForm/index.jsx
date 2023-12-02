import React, {lazy} from 'react';
import api from '../../api';

const Form = lazy(() => import("../../components/Form/index.jsx"));
const FormField = lazy(() => import("../../components/FormField/index.jsx"));
const FormButton = lazy(() => import("../../components/FormButton/index.jsx"));

const sendRegisterRequest = async e => {
  e.preventDefault();

  let name = e.target.elements['name'].value;
  let password = e.target.elements['password'].value;
  if (!name) {console.error("name required"); return;}
  if (!password) {console.error("password required"); return;}

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
    <Form name="register-form" onSubmit={sendRegisterRequest}>
      <FormField required name="name" type="text" />
      <FormField required name="password" type="password" />
      <FormButton type="submit" text="Register" />
    </Form>
  );
}

export default RegisterForm;
