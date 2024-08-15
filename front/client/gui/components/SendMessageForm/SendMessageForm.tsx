import React, {lazy, useCallback, useState, useRef} from 'react';
import api from '../../api';
import Form from '../Form';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

const SendMessageForm = ({ onError, onSuccess }) => {
  const fileRef = useRef(null);

  const [message, setMessage] = useState("");
  const [file, setFile] = useState(null);

  const onFileChange = useCallback(e => {
    setFile(e.target.files[0])
  }, [setFile]);

  const onMessageChange = useCallback(e => {
    setMessage(e.target.value);
  }, [setMessage]);


  const sendMessageRequest = useCallback(e => {
    e.preventDefault();

    console.log("send message request");

    const send = async form => {
      try {
        const response = await api("/messages/v1/send", {
          method: "POST",
          credentials: "include",
          body: form,
        });
        onSuccess(response)

        fileRef.current.value = null
        setMessage("")
        setFile(null)
      } catch (e) {
        console.error(i18n("error_occured"), e);
        onError(e);
      }
    }

    if (!message) {console.error(i18n("msg_required_err")); return;}

    const form = new FormData()
    form.append("message", message);
    if (file.name != "") form.append('file', file, file.name);

    send(form)
  }, [
    onSuccess,
    onError,
    setMessage,
    setFile,
    message,
    file,
    fileRef,
  ]);

  return (
    <>
      <Form
        autoComplete="off"
        css="row items-end"
      >
        <FormField
          required
          el="textarea"
          name="message"
          type="input"
          value={message}
          onChange={onMessageChange}
        />
        <FormField
          ref={fileRef}
          el="input"
          name="file"
          type="file"
          onChange={onFileChange}
        />
      </Form>

      <Button
        type="button"
        text={i18n("msg_send_text")}
        onClick={sendMessageRequest}
      />
    </>
  );
}

export default SendMessageForm;
