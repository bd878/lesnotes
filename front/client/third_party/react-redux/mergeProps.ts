function defaultMergeProps(stateProps, dispatchProps, ownProps) {
  return { ...ownProps, ...stateProps, ...dispatchProps }
}

export function mergePropsFactory(_mergeProps) {
  return () => defaultMergeProps
}