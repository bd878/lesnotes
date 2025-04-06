import React from 'react'
import {createSubscription} from './Subscription'
import {useIsomorphicLayoutEffect} from './useIsomorphicLayoutEffect'
import {ReactReduxContext} from './Context'

export default function Provider(providerProps) {
  const { children, context, store } = providerProps

  const contextValue = React.useMemo(() => {
    const subscription = createSubscription(store)

    const baseContextValue = {
      store,
      subscription,
    }

    return baseContextValue
  }, [store])

  const previousState = React.useMemo(() => store.getState(), [store])

  useIsomorphicLayoutEffect(() => {
    const { subscription } = contextValue
    subscription.onStateChange = subscription.notifyNestedSubs
    subscription.trySubscribe()

    // compare by pointer
    if (previousState !== store.getState()) {
      subscription.notifyNestedSubs()
    }
    return () => {
      subscription.tryUnsubscribe()
      subscription.onStateChange = undefined
    }
  }, [contextValue, previousState])

  const Context = context || ReactReduxContext

  return <Context.Provider value={contextValue}>{children}</Context.Provider>
}