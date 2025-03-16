import React, {useCallback, useMemo, useEffect, useState, useRef} from 'react';
import Form from '../Form';
import i18n from '../../i18n';
import MessageFormComponent from './MessageFormComponent';

function MessageFormContainer(props) {
	const {
		send,
		update,
		messageForEdit,
		reset,
	} = props

	const fileRef = useRef(null);

	const [message, setMessage] = useState("");
	const [file, setFile] = useState(null);

	const isEditMode = useMemo(() => messageForEdit && messageForEdit.ID, [messageForEdit])

	useEffect(() => {
		if (isEditMode)
			setMessage(messageForEdit.text)
	}, [setMessage, isEditMode]);

	const onFileChange = useCallback(e => {
		setFile(e.target.files[0])
	}, [setFile]);

	const onMessageChange = useCallback(e => {
		setMessage(e.target.value);
	}, [setMessage]);

	const updateMessageRequest = useCallback(e => {
		update(messageForEdit.ID, message)
		setMessage("")
		setFile(null)
	}, [update, setMessage, setFile, messageForEdit, message])

	const sendMessageRequest = useCallback(e => {
		if (!message) {console.error(i18n("msg_required_err")); return;}

		send(message, file)

		fileRef.current.value = null
		setMessage("")
		setFile(null)
	}, [send, setMessage, setFile, message, file, fileRef]);

	const onSubmit = useCallback(e => {
		e.preventDefault()
		if (messageForEdit && messageForEdit.ID) 
			updateMessageRequest(e) // update mode
		else
			sendMessageRequest(e) // save mode
	}, [sendMessageRequest, updateMessageRequest, message, messageForEdit])

	const onEditCancel = useCallback(() => {
		reset()
		setMessage("")
	}, [reset, setMessage])

	return (
		<MessageFormComponent
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

export default MessageFormContainer
