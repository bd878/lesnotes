function bindActionCreator(
  actionCreator,
  dispatch
) {
  return function(this: any, ...args: any[]) {
    return dispatch(actionCreator.apply(this, args))
  }
}

export default function bindActionCreators(
  actionCreators,
  dispatch
) {
  if (typeof actionCreators === "function") {
    return bindActionCreators(actionCreators, dispatch)
  }

  if (typeof actionCreators !== "object" || actionCreators === null) {
    throw new Error("actionCreators must be a non-null object")
  }

  const boundActionCreators = {}
  for (const key in actionCreators) {
    const actionCreator = actionCreators[key]
    if (typeof actionCreator === "function") {
      boundActionCreators[key] = bindActionCreator(actionCreator, dispatch)
    }
  }
  return boundActionCreators
}