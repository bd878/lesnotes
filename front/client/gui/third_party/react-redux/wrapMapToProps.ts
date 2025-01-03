export function wrapMapToPropsConstant(getConstant) {
  return function initConstantSelector(dispatch) {
    const constant = getConstant(dispatch)
    function constantSelector() {
      return constant
    }
    constantSelector.dependsOnOwnProps = false
    return constantSelector
  }
}

export function wrapMapToPropsFunc(realMapToProps, methodName) {
  return function initProxySelector(
    dispatch,
    { displayName }
  ) {
    const proxy = function mapToPropsProxy(
      stateOrDispatch,
      ownProps,
    ) {
      return proxy.mapToProps(stateOrDispatch, ownProps)
    }

    proxy.mapToProps = function detectFactoryAndVerify(
      stateOrDispatch,
      ownProps,
    ) {
      proxy.mapToProps = realMapToProps
      let props = proxy(stateOrDispatch, ownProps)

      if (typeof props === "function") {
        proxy.mapToProps = props
        props = proxy(stateOrDispatch, ownProps)
      }

      return props
    }

    return proxy
  }
}