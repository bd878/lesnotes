import React, {useCallback, useState, useRef} from 'react';
import {connect} from '../../third_party/react-redux';
import Form from '../Form';
import i18n from '../../i18n';
import SendMessageFormComponent from './SendMessageFormComponent';
import {sendMessageActionCreator} from '../../features/messages';

function SendMessageFormContainer(props) {
  const {
    onSuccess,
    onError,
    sendMessage,
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

    if (!message) {console.error(i18n("msg_required_err")); return;}

    console.log("message=", message, "file=", file)
    sendMessage(message, file)

    fileRef.current.value = null
    setMessage("")
    setFile(null)
  }, [sendMessage, setMessage, setFile,
    message, file, fileRef]);

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
  sendMessage: sendMessageActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(
  SendMessageFormContainer
)