import React, {lazy, useCallback} from 'react';
import api from '../../api';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form/Form.tsx"));
const FormField = lazy(() => import("../../components/FormField/FormField.tsx"));
const Button = lazy(() => import("../../components/Button/Button.tsx"));

const SendMessageForm = ({ setError, onSend }) => {
  const sendMessageRequest = useCallback(e => {
    e.preventDefault();

    const send = async form => {
      try {
        await api("/messages/v1/send", {
          method: "POST",
          credentials: "include",
          body: form,
        });
      } catch (e) {
        console.error(i18n("error_occured"), e);
        setError(i18n("loading_messages_error"));
        return
      }

      onSend()
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
  }, [setError, onSend]);

  return (
    <Form
      name="send-message-form"
      onSubmit={sendMessageRequest}
      enctype="multipart/form-data"
    >
      <FormField required name="message" type="text" />
      <FormField name="file" type="file" />
      <Button type="submit" text={i18n("msg_send_text")} />
    </Form>
  );
}

export default SendMessageForm;
