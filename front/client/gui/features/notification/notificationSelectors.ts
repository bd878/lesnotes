export const selectIsNotificationVisible = state => state.notification.isVisible || false
export const selectNotificationText = state => state.notification.text || ""
export const selectTimerID = state => state.notification.timerID || null