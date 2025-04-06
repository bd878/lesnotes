const ActionTypes = {INIT: "init", REPLACE: "replace"}

const callListener = listener => listener();

export function createStore(reducer, preloadedState, enhancer) {
  if (typeof enhancer !== "undefined") {
    return enhancer(createStore)(
      reducer,
      preloadedState
    )
  }

  let currentReducer = reducer
  let currentState = preloadedState || undefined
  // window
  let currentListeners = new Map()
  let nextListeners = currentListeners 
  let listenerIdCounter = 0
  let isDispatching = false

  function getState() {
    if (isDispatching) {
      throw new Error("cannot call store.getState() while reducer is executing")
    }

    return currentState
  }

  /**
   * deep copy listeners.
   * Ensure that listeners are not mutated,
   * if dispatch is in progress.
   * We update current listeners with next listeners
   * just before listeners are notified in
   * each dispatch
   */
  function ensureCanMutateNextListeners() {
    if (nextListeners === currentListeners) {
      nextListeners = new Map()
      currentListeners.forEach((listener, key) => {
        nextListeners.set(key, listener)
      })
    }
  }

  function subscribe(listener) {
    if (typeof listener !== 'function') {
      throw new Error("listener is not a function")
    }

    if (isDispatching) {
      throw new Error("cannot call store.subscribe() while reducer is executing")
    }

    let isSubscribed = true
    ensureCanMutateNextListeners()
    const listenerId = listenerIdCounter++
    nextListeners.set(listenerId, listener)

    return function unsibscribe() {
      if (!isSubscribed) {
        return
      }

      if (isDispatching) {
        throw new Error("may not unsubscribe while reducer is executing")
      }

      isSubscribed = false

      ensureCanMutateNextListeners()
      nextListeners.delete(listenerId)
      currentListeners = null
    }
  }

  function dispatch(action) {
    if (typeof action !== "object") {
      throw new Error("actions must be object")
    }

    if (typeof action.type === "undefined") {
      throw new Error("actions must have \"type\" property")
    }

    if (typeof action.type !== "string") {
      throw new Error("action property \"type\" must be a string")
    }

    if (isDispatching) {
      throw new Error("reducers may not dispatch actions")
    }

    try {
      isDispatching = true
      currentState = currentReducer(currentState, action)
    } catch(e) {
      console.log(e)
    } finally {
      isDispatching = false
    }

    currentListeners = nextListeners;
    currentListeners.forEach(callListener)
    return action;
  }

  function replaceReducer(nextReducer) {
    if (typeof nextReducer !== 'function') {
      throw new Error("next reducer must be a function")
    }

    currentReducer = nextReducer
    dispatch({ type: ActionTypes.REPLACE })
  }

  dispatch({ type: ActionTypes.INIT })

  const store = {
    dispatch,
    subscribe,
    getState,
    replaceReducer
  }

  return store;
}
