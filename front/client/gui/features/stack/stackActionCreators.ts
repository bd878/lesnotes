import {
	OPEN_THREAD,
	CLOSE_THREAD,
	DESTROY_THREAD,
	PUBLISH_SELECTED,
	PRIVATE_SELECTED,

	UPDATE_MESSAGE,
	DELETE_MESSAGE,
	DELETE_SELECTED,
	SEND_MESSAGE,
	COPY_MESSAGE,
	COPY_LINK,
	FETCH_MESSAGES,
	SELECT_MESSAGE,
	UNSELECT_MESSAGE,
	CLEAR_SELECTED,
	SET_MESSAGE_FOR_EDIT,

	MESSAGES_FAILED,

	SEND_MESSAGE_SUCCEEDED,
	FETCH_MESSAGES_SUCCEEDED,
	UPDATE_MESSAGE_SUCCEEDED,
	DELETE_MESSAGE_SUCCEEDED,
	DELETE_SELECTED_SUCCEEDED,
} from './stackActions';

export const openThreadActionCreator = payload => ({
	type: OPEN_THREAD,
	payload,
})

export const closeThreadActionCreator = payload => ({
	type: CLOSE_THREAD,
	payload,
})

export const destroyThreadActionCreator = payload => ({
	type: DESTROY_THREAD,
	payload,
})

export const messagesFailedActionCreator = index => payload => ({
	type: MESSAGES_FAILED,
	index: index,
	payload,
})

export const sendMessageActionCreator = index => payload => ({
	type: SEND_MESSAGE,
	index: index,
	payload,
})

export const sendMessageSucceededActionCreator = index => payload => ({
	type: SEND_MESSAGE_SUCCEEDED,
	index: index,
	payload,
})

export const fetchMessagesActionCreator = index => payload => ({
	type: FETCH_MESSAGES,
	index: index,
	payload,
})

export const fetchMessagesSucceededActionCreator = index => payload => ({
	type: FETCH_MESSAGES_SUCCEEDED,
	index: index,
	payload,
})

export const updateMessageActionCreator = index => payload => ({
	type: UPDATE_MESSAGE,
	index,
	payload,
})

export const updateMessageSucceededActionCreator = index => payload => ({
	type: UPDATE_MESSAGE_SUCCEEDED,
	index: index,
	payload,
})

export const setEditMessageActionCreator = index => payload => ({
	type: SET_MESSAGE_FOR_EDIT,
	index: index,
	payload,
})

export const resetEditMessageActionCreator = index => setEditMessageActionCreator(index)

export const deleteMessageActionCreator = index => payload => ({
	type: DELETE_MESSAGE,
	index: index,
	payload,
})

export const deleteSelectedActionCreator = index => () => ({
	type: DELETE_SELECTED,
	index: index,
})

export const selectMessageActionCreator = index => payload => ({
	type: SELECT_MESSAGE,
	index: index,
	payload,
})

export const deleteSelectedSucceededActionCreator = index => payload => ({
	type: DELETE_SELECTED_SUCCEEDED,
	index: index,
	payload: payload,
})

export const copyLinkActionCreator = index => payload => ({
	type: COPY_LINK,
	index: index,
	payload: payload,
})

export const unselectMessageActionCreator = index => payload => ({
	type: UNSELECT_MESSAGE,
	index: index,
	payload,
})

export const clearSelectedActionCreator = index => () => ({
	type: CLEAR_SELECTED,
	index: index,
})

export const copyMessageActionCreator = index => payload => ({
	type: COPY_MESSAGE,
	index: index,
	payload,
})

export const deleteMessageSucceededActionCreator = index => payload => ({
	type: DELETE_MESSAGE_SUCCEEDED,
	index: index,
	payload,
})

export const publishSelectedActionCreator = index => () => ({
	type: PUBLISH_SELECTED,
	index: index,
})

export const privateSelectedActionCreator = index => () => ({
	type: PRIVATE_SELECTED,
	index: index,
})