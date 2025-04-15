export const selectStack = state => state.stack.stack || []

export const selectMessages = index => state => state.stack.stack[index].list || []
export const selectIsLastPage = index => state => state.stack.stack[index].isLastPage || false
export const selectIsLoading = index => state => state.stack.stack[index].loading || false
export const selectError = index => state => state.stack.stack[index].error || ""
export const selectLoadOffset = index => state => selectMessages(index)(state).length
export const selectMessageForEdit = index => state => state.stack.stack[index].messageForEdit
export const selectIsEditMode = index => state => selectMessageForEdit(index)(state).ID ? true : false
