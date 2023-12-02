import React, {lazy} from 'react';
import api from '../../api';

const Form = lazy(() => import("../../components/Form/index.jsx"));
const FormField = lazy(() => import("../../components/FormField/index.jsx"));
const FormButton = lazy(() => import("../../components/FormButton/index.jsx"));

const sendMessageRequest = async e => {
  e.preventDefault();

  let message = e.target.elements['message'].value;
  if (!message) {console.error("message required"); return;}

  const form = new FormData()
  form.append("message", message);
  const response = await api("/messages/v1/send", {
    method: "POST",
    credentials: "include",
    body: form,
  });
  console.log(response);
}

const SendMessageForm = props => {
  return (
    <Form name="send-message-form" onSubmit={sendMessageRequest}>
      <FormField required name="message" type="text" />
      <FormButton type="submit" text="Send" />
    </Form>
  );
}

export default SendMessageForm;
