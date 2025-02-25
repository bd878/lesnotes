import React, {useCallback, useEffect, useState, useRef} from 'react';
import {connect} from '../../third_party/react-redux';
import Form from '../Form';
import i18n from '../../i18n';
import SendMessageFormComponent from './SendMessageFormComponent';
import {
  sendMessageActionCreator,
  updateMessageActionCreator,
  selectMessageForEdit,
  selectIsEditMode,
  resetEditMessageActionCreator,
} from '../../features/messages';

function SendMessageFormContainer(props) {
  const {
    onSuccess,
    onError,
    sendMessage,
    updateMessage,
    messageForEdit,
    isEditMode,
    resetEditMessage,
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
    setMessage("")
    setFile(null)
  }, [updateMessage, setMessage,
    setFile, messageForEdit, message])

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

  const onEditCancel = useCallback(() => {
    resetEditMessage()
    setMessage("")
  }, [resetEditMessage, setMessage])

  return (
    <SendMessageFormComponent
      fileRef={fileRef}
      message={message}
      onMessageChange={onMessageChange}
      onFileChange={onFileChange}
      onSubmit={onSubmit}
      shouldShowCancelButton={isEditMode}
      onCancel={onEditCancel}
    />
  )
}

const mapStateToProps = state => ({
  messageForEdit: selectMessageForEdit(state),
  isEditMode: selectIsEditMode(state),
})

const mapDispatchToProps = ({
  sendMessage: sendMessageActionCreator,
  updateMessage: updateMessageActionCreator,
  resetEditMessage: resetEditMessageActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(
  SendMessageFormContainer
)