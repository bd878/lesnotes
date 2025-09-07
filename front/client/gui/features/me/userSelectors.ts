import models from '../../../api/models'

export const selectUser = state => state.me.user || models.user()
export const selectIsDesktop = state => state.me.isDesktop || false
export const selectIsMobile = state => state.me.isMobile || false
export const selectIsMiniapp = state => state.me.isMiniapp || false
export const selectBrowser = state => state.me.browser || ""
