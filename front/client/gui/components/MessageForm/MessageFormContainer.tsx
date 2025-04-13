import React, {useCallback, useMemo, useEffect, useState, useRef} from 'react';
import Form from '../Form';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is';
import MessageFormComponent from './MessageFormComponent';

function MessageFormContainer(props) {
	const {
		send,
		update,
		messageForEdit,
		reset,
	} = props

	const fileRef = useRef(null);

	const [text, setText] = useState("");
	const [file, setFile] = useState(null);

	const isEditMode = useMemo(() => is.notEmpty(messageForEdit) && is.notEmpty(messageForEdit.text), [messageForEdit])

	useEffect(() => {
		if (isEditMode)
			setText(messageForEdit.text)
	}, [setText, isEditMode]);

	const onFileChange = useCallback(e => {
		setFile(e.target.files[0])
	}, [setFile]);

	const onMessageChange = useCallback(e => {
		setText(e.target.value);
	}, [setText]);

	const updateMessageRequest = useCallback(e => {
		update(messageForEdit.ID, text)
		setText("")
		setFile(null)
	}, [update, setText, setFile, messageForEdit, text])

	const sendMessageRequest = useCallback(e => {
		if (!text) {console.error(i18n("msg_required_err")); return;}

		send(text, file)

		fileRef.current.value = null
		setText("")
		setFile(null)
	}, [send, setText, setFile, text, file, fileRef]);

	const onSubmit = useCallback(e => {
		e.preventDefault()
		if (messageForEdit && messageForEdit.ID) 
			updateMessageRequest(e) // update mode
		else
			sendMessageRequest(e) // save mode
	}, [sendMessageRequest, updateMessageRequest, text, messageForEdit])

	const onEditCancel = useCallback(() => {
		reset()
		setText("")
	}, [reset, setText])

	const onFileClick = useCallback((e) => {
		e.preventDefault()
		if (is.notEmpty(fileRef))
			fileRef.current.click()
	}, [fileRef])

	return (
		<MessageFormComponent
			fileRef={fileRef}
			text={text}
			file={file}
			onFileClick={onFileClick}
			onMessageChange={onMessageChange}
			onFileChange={onFileChange}
			onSubmit={onSubmit}
			shouldShowCancelButton={isEditMode}
			onCancel={onEditCancel}
		/>
	)
}

export default MessageFormContainer
