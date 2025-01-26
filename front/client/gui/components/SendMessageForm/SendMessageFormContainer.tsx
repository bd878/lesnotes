import React, {lazy, useCallback, useState, useRef} from 'react';
import {connect} from '../../third_party/react-redux';
import api from '../../api';
import Form from '../Form';
import i18n from '../../i18n';
import SendMessageFormComponent from './SendMessageFormComponent';
import {appendMessagesActionCreator} from '../../features/messages';

function SendMessageFormContainer(props) {
  const {
    onSuccess,
    onError,
    appendMessage,
  } = props

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

    const send = async () => {
      try {
        const response = await api.sendMessage(message, file)
        if (response.error != "") {
          console.error(i18n("error_occured"), response.error)
          console.log(response.explain)
          return
        }
        fileRef.current.value = null
        setMessage("")
        setFile(null)

        appendMessage([response.message])
      } catch (e) {
        console.error(i18n("error_occured"), e)
        return
      }
    }

    if (!message) {console.error(i18n("msg_required_err")); return;}

    send()
  }, [
    appendMessage,
    onSuccess,
    onError,
    setMessage,
    setFile,
    message,
    file,
    fileRef,
  ]);

  return (
    <SendMessageFormComponent
      fileRef={fileRef}
      message={message}
      onMessageChange={onMessageChange}
      onFileChange={onFileChange}
      sendMessageRequest={sendMessageRequest}
    />
  )
}

const mapStateToProps = () => ({})

const mapDispatchToProps = ({
  appendMessage: appendMessagesActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(
  SendMessageFormContainer
)