export const selectMessages = state => state.messages.list || []
export const selectIsLastPage = state => state.messages.isLastPage || false
export const selectIsLoading = state => state.messages.loading || false
export const selectError = state => state.messages.error || ""
export const selectLoadOffset = state => selectMessages(state).length
export const selectMessageForEdit = state => state.messages.messageForEdit
export const selectIsEditMode = state => selectMessageForEdit(state).ID ? true : false
