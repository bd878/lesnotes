function defaultMergeProps(stateProps, dispatchProps, ownProps) {
  return { ...ownProps, ...stateProps, ...dispatchProps }
}

export function mergePropsFactory(mergeProps) {
  return !mergeProps
    ? () => defaultMergeProps
    : throw new Error("mergeProps is not supported yet")
}