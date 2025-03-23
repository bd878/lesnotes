import React, {useState, useEffect, useCallback} from 'react';
import api from '../../api';
import models from '../../api/models';
import * as is from '../../third_party/is'
import MessagePageComponent from './MessagePageComponent';

function MessagePageContainer(props) {
	const {
		id,
	} = props

	const [message, setMessage] = useState(models.message())
	const [messageForEdit, setMessageForEdit] = useState(models.message())
	const [error, setError] = useState(false)
	const [loading, setLoading] = useState(false)

	useEffect(() => {
		async function loadMessage(id) {
			setLoading(true)

			const result = await api.loadOneMessage(id)
			if (is.notEmpty(result.error)) {
				console.error(result.error, result.explain)
				setError(result.error)
				return
			}
			setMessage(result.message)

			setLoading(false)
		}

		if (is.notEmpty(id))
			loadMessage(id)
	}, [id, setError, setMessage, setLoading])

	const sendMessage = useCallback((text, file) => {
		async function sendMessage(text, file) {
			setError(false)
			setLoading(true)

			const result = await api.sendMessage({text, file})
			if (is.notEmpty(result.error)) {
				console.error(result.error, result.explain)
				setError(result.error)
				return
			}
			setMessage(result.message)
			setLoading(false)
		}

		sendMessage(text, file)
	}, [setLoading, setError, setMessage])

	const updateMessage = useCallback(() => {}, [])

	const resetEditMessage = useCallback(() => {}, [])

	return (
		<MessagePageComponent
			message={message}
			error={error}
			loading={loading}
			resetEditMessage={resetEditMessage}
			updateMessage={updateMessage}
			sendMessage={sendMessage}
			messageForEdit={messageForEdit}
			shouldRenderEditForm={is.empty(id)}
		/>
	)
}

export default MessagePageContainer