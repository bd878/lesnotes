import React, {lazy} from 'react';
import api from '../../api';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form/index.jsx"));
const FormField = lazy(() => import("../../components/FormField/index.jsx"));
const Button = lazy(() => import("../../components/Button/index.jsx"));

const sendMessageRequest = async e => {
  e.preventDefault();

  let message = e.target.elements['message'].value;
  if (!message) {console.error(i18n("msg_required_err")); return;}
  let file = e.target.elements['file'].value;

  const form = new FormData()
  form.append("message", message);
  if (file) {form.append('file', file);}

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
      <FormField name="file" type="file" />
      <Button type="submit" text={i18n("msg_send_text")} />
    </Form>
  );
}

export default SendMessageForm;
