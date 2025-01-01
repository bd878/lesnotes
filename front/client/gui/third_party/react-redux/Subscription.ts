function batch(callback) {
  callback()
}

function createListenersCollection() {
  let first = null
  let last = null

  return {
    clear() {
      first = null
      last = null
    },

    notify() {
      batch(() => {
        let listener = first
        while (listener) {
          listener.callback()
          listener = listener.next
        }
      })
    },

    get() {
      const listeners = []
      let listener = first
      while (listener) {
        listeners.push(listener)
        listener = listener.next
      }
      return listeners
    },

    subscribe(callback) {
      let isSubscribed = true

      const listener = (last = {
        callback,
        next: null,
        prev: last,
      })

      if (listener.prev) {
        listener.prev.next = listener
      } else {
        first = listener
      }

      return function unsubscribe() {
        if (!isSubscribed || first === null) return
        isSubscribed = false

        if (listener.next) {
          listener.next.prev = listener.prev
        } else {
          last = listener.prev
        }

        if (listener.prev) {
          listener.prev.next = listener.next
        } else {
          first = listener.next
        }
      }
    },
  }
}

const nullListeners = {
  notify() {},
  get: () => [],
}

export function createSubscription(store, parentSub) {
  let unsubscribe = undefined
  let listeners = nullListeners

  let subscriptionsAmount = 0
  let selfSubscribed = false

  function addNestedSub(listener) {
    trySubscribe()

    const cleanupListener = listeners.subscribe(listener)

    let removed = false
    return () => {
      if (!removed) {
        removed = true
        cleanupListener()
        tryUnsubscribe()
      }
    }
  }

  function notifyNestedSubs() {
    listeners.notify()
  }

  function handleChangeWrapper() {
    if (subscription.onStateChange) {
      subscription.onStateChange()
    }
  }

  function isSubscribed() {
    return selfSubscribed
  }

  function tryUnsubscribe() {
    subscriptionsAmount++
    if (!unsubscribe) {
      unsubscribe = parentSub
        ? parentSub.addNestedSub(handleChangeWrapper)
        : store.subscribe(handleChangeWrapper)

      listeners = createListenersCollection()
    }
  }

  function tryUnsubscribe() {
    subscriptionsAmount--
    if (unsubscribe && subscriptionsAmount === 0) {
      unsubscribe()
      unsubscribe = undefined
      listeners.clear()
      listeners = nullListeners
    }
  }

  function trySubscribeSelf() {
    if (!selfSubscribed) {
      selfSubscribed = true
      trySubscribe()
    }
  }

  function tryUnsibscribeSelf() {
    if (selfSubscribed) {
      selfSubscribed = false
      tryUnsubscribe()
    }
  }

  const subscription = {
    addNestedSub,
    notifyNestedSubs,
    isSubscribed,
    trySubscribe: trySubscribeSelf,
    tryUnsubscribe: tryUnsibscribeSelf,
    getListeners: () => listeners,
  }

  return subscription
}