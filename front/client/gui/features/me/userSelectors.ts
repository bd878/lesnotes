import models from '../../../api/models'

export const selectIsLoading = state => state.me.loading || false
export const selectError = state => state.me.error || ""
export const selectIsError = state => selectError(state) || false
export const selectWillRedirect = state => state.me.willRedirect || false
export const selectUser = state => state.me.user || models.user()
export const selectIsAuth = state => state.me.isAuth || false
export const selectIsDesktop = state => state.me.isDesktop || false
export const selectIsMobile = state => state.me.isMobile || false
export const selectIsMiniapp = state => state.me.isMiniapp || false
export const selectBrowser = state => state.me.browser || ""
