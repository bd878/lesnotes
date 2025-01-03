import {wrapMapToPropsConstant, wrapMapToPropsFunc} from './wrapMapToProps'
import bindActionCreators from './bindActionCreators'

export function mapDispatchToPropsFactory(mapDispatchToProps) {
  return mapDispatchToProps && typeof mapDispatchToProps === 'object'
    ? wrapMapToPropsConstant(dispatch =>
        bindActionCreators(mapDispatchToProps, dispatch)
      )
    : !mapDispatchToProps
      ? wrapMapToPropsConstant(dispatch => ({dispatch}))
      : wrapMapToPropsFunc(mapDispatchToProps, 'mapDispatchToProps')
}