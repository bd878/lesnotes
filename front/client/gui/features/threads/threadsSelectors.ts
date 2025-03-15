export const selectMessages = state => state.threads.list || []
export const selectIsLastPage = state => state.threads.isLastPage || false
export const selectIsLoading = state => state.threads.loading || false
export const selectError = state => state.threads.error || ""
export const selectLoadOffset = state => selectMessages(state).length
export const selectThreadMessage = state => state.threads.message
export const selectThreadID = state => state.threads.threadID || 0