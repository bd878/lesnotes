import React, {lazy, useCallback} from 'react';
import api from '../../api';
import Form from '../Form';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

const SendMessageForm = ({ onError, onSuccess }) => {
  const sendMessageRequest = useCallback(e => {
    e.preventDefault();

    const send = async form => {
      try {
        const response = await api("/messages/v1/send", {
          method: "POST",
          credentials: "include",
          body: form,
        });
        onSuccess(response)
      } catch (e) {
        console.error(i18n("error_occured"), e);
        onError(e);
      }
    }

    let message = e.target.elements['message'].value;
    if (!message) {console.error(i18n("msg_required_err")); return;}
    let file = e.target.elements['file'].files[0];
    let filename = e.target.elements['file'].value;

    const form = new FormData()
    form.append("message", message);
    if (file) {
      form.append('file', file, filename);
    }

    send(form)
  }, [onSuccess, onError]);

  return (
    <Form
      name="send-message-form"
      onSubmit={sendMessageRequest}
      enctype="multipart/form-data"
    >
      <FormField required name="message" type="textarea" />
      <FormField name="file" type="file" />
      <Button type="submit" text={i18n("msg_send_text")} />
    </Form>
  );
}

export default SendMessageForm;
