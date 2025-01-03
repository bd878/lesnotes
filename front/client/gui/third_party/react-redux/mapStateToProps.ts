import {wrapMapToPropsConstant, wrapMapToPropsFunc} from './wrapMapToProps'

export function mapStateToPropsFactory(mapStateToProps) {
  return !mapStateToProps
    ? wrapMapToPropsConstant(() => ({}))
    : wrapMapToPropsFunc(mapStateToProps, 'mapStateToProps')
}