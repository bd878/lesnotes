export default function combineReducers(reducers) {
  /**
   * ensure that reducers is a key:function map
   * and all keys are defined
   */
  const reducerKeys = Object.keys(reducers)
  const finalReducers = {}
  for (let i = 0; i < reducerKeys.length; i++) {
    const key = reducerKeys[i]

    if (typeof reducers[key] === "function") {
      finalReducers[key] = reducers[key]
    }
  }
  const finalreducerKeys = Object.keys(finalReducers)

  /* TODO: reducers shape assertion */

  return function combination(
    state,
    action
  ) {
    let hasChanged = false
    const nextState = {}
    for (let i = 0; i < finalreducerKeys.length; i++) {
      const key = finalreducerKeys[i]
      const reducer = finalReducers[key]
      const previousStateForKey = state[key]
      const nextStateForKey = reducer(previousStateForKey, action)

      if (typeof nextStateForKey === "undefined") {
        throw new Error(`reducer for action type ${action.type} returned undefined`)
      }
      nextState[key] = nextStateForKey
      hasChanged = hasChanged || nextStateForKey !== previousStateForKey
    }
    hasChanged = hasChanged || finalreducerKeys.length !== Object.keys(state).length
    return hasChanged ? nextState : state
  }
}