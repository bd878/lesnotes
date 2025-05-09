export const selectStack = state => state.stack.stack || []

export const selectSelectedMessageIDs = index => state => state.stack.stack[index].selectedMessageIDs || new Set()
export const selectHasNextThread = index => state => state.stack.stack.length > index+1
export const selectThreadID = index => state => state.stack.stack[index].ID || 0
export const selectMessages = index => state => state.stack.stack[index].list || []
export const selectIsLastPage = index => state => state.stack.stack[index].isLastPage || false
export const selectIsLoading = index => state => state.stack.stack[index].loading || false
export const selectError = index => state => state.stack.stack[index].error || ""
export const selectLoadOffset = index => state => selectMessages(index)(state).length
export const selectMessageForEdit = index => state => state.stack.stack[index].messageForEdit
export const selectIsEditMode = index => state => selectMessageForEdit(index)(state).ID ? true : false
export const selectIsMessageThreadOpen = index => state => messageID => selectHasNextThread(index)(state) ? state.stack.stack[index+1].ID === messageID : false
