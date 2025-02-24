import React, {useCallback, useEffect, useState, useRef} from 'react';
import {connect} from '../../third_party/react-redux';
import Form from '../Form';
import i18n from '../../i18n';
import SendMessageFormComponent from './SendMessageFormComponent';
import {
  sendMessageActionCreator,
  updateMessageActionCreator,
  selectMessageForEdit,
} from '../../features/messages';

function SendMessageFormContainer(props) {
  const {
    onSuccess,
    onError,
    sendMessage,
    updateMessage,
    messageForEdit,
  } = props

  const fileRef = useRef(null);

  const [message, setMessage] = useState("");
  const [file, setFile] = useState(null);

  useEffect(() => {
    if (messageForEdit && messageForEdit.ID)
      setMessage(messageForEdit.text)
  }, [setMessage, messageForEdit]);

  const onFileChange = useCallback(e => {
    setFile(e.target.files[0])
  }, [setFile]);

  const onMessageChange = useCallback(e => {
    setMessage(e.target.value);
  }, [setMessage]);

  const updateMessageRequest = useCallback(e => {
    updateMessage(messageForEdit.ID, message)
  }, [updateMessage, setMessage, messageForEdit, message])

  const sendMessageRequest = useCallback(e => {
    if (!message) {console.error(i18n("msg_required_err")); return;}

    sendMessage(message, file)

    fileRef.current.value = null
    setMessage("")
    setFile(null)
  }, [sendMessage, setMessage, setFile,
    message, file, fileRef]);

  const onSubmit = useCallback(e => {
    e.preventDefault()
    if (messageForEdit && messageForEdit.ID)
      updateMessageRequest(e) // update mode
    else
      sendMessageRequest(e) // save mode
  }, [sendMessageRequest, updateMessageRequest,
      message, messageForEdit])

  return (
    <SendMessageFormComponent
      fileRef={fileRef}
      message={message}
      onMessageChange={onMessageChange}
      onFileChange={onFileChange}
      onSubmit={onSubmit}
    />
  )
}

const mapStateToProps = state => ({
  messageForEdit: selectMessageForEdit(state),
})

const mapDispatchToProps = ({
  sendMessage: sendMessageActionCreator,
  updateMessage: updateMessageActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(
  SendMessageFormContainer
)