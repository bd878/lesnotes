import React from 'react';
import hoistStatics from './hoistStatics';
import shallowEqual from './shallowEqual';
import {mapStateToPropsFactory} from './mapStateToProps'
import {mapDispatchToPropsFactory} from './mapDispatchToProps'
import {mergePropsFactory} from './mergeProps'
import {createSubscription} from './Subscription'
import {ReactReduxContext} from './Context'
import {useIsomorphicLayoutEffect} from './useIsomorphicLayoutEffect'

function strictEqual(a, b) {
  return a === b
}

function captureWrapperProps(
  lastProps,
  lastChildProps,
  renderIsScheduled,
  props,
  childPropsFromStoreUpdate,
  notifyNestedSubs
) {
  lastProps.current = props
  renderIsScheduled.current = false

  if (childPropsFromStoreUpdate.current) {
    childPropsFromStoreUpdate.current = null
    notifyNestedSubs()
  }
}

function subscribeUpdates(
  store,
  subscription,
  childPropsSelector,
  lastProps,
  lastChildProps,
  renderIsScheduled,
  isMounted,
  childPropsFromStoreUpdate,
  notifyNestedSubs,
  listener
) {
  let didUnsubscribe = false
  let lastThrownError = null

  const checkForUpdates = () => {
    if (didUnsubscribe || !isMounted.current) {
      return
    }

    const latestStoreState = store.getState()

    let newChildProps, error
    try {
      newChildProps = childPropsSelector(
        latestStoreState,
        lastProps.current,
      )
    } catch(e) {
      error = e
      lastThrownError = e
    }

    if (!error) {
      lastThrownError = null
    }

    if (newChildProps === lastChildProps.current) {
      if (!renderIsScheduled.current) {
        notifyNestedSubs()
      }
    } else {
      // save references to the new child props
      lastChildProps.current = newChildProps
      childPropsFromStoreUpdate.current = newChildProps
      renderIsScheduled.current = true

      // trigger React listener
      listener()
    }
  }

  subscription.onStateChange = checkForUpdates
  subscription.trySubscribe()

  checkForUpdates()

  const unsubscribeWrapper = () => {
    didUnsubscribe = true
    subscription.tryUnsubscribe()
    subscription.onStateChange = null

    if (lastThrownError) {
      throw lastThrownError
    }
  }

  return unsubscribeWrapper
}

export function connect(mapStateToProps, mapDispatchToProps) {
  const initMapStateToProps = mapStateToPropsFactory(mapStateToProps)
  const initMapDispatchToProps = mapDispatchToPropsFactory(mapDispatchToProps)
  const initMergeProps = mergePropsFactory(mergeProps)

  const wrapWithConnect = (WrappedComponent) => {
    const wrappedComponentName =
      WrappedComponent.displayName || WrappedComponent.name || 'Component';

    const displayName = `Connect(${wrappedComponentName})`

    // selector factory merges mapStateToProps and mapDispatchToProps props
    // in one component props
    const selectorFactoryOptions = {
      displayName,
      WrappedComponent,
      initMapStateToProps,
      initMapDispatchToProps,
      initMergeProps,
      areStatesEqual: strictEqual,
      areOwnPropsEqual: shallowEqual,
      areStatePropsEqual: shallowEqual,
    }

    function ConnectFunction(props) {
      const contextValue = React.useContext(ReactReduxContext)

      const store = contextValue.store

      // mapStateToProps, mapDispatchToProps, mergeProps : depends on store
      const childPropsSelector = React.useMemo(() => {
        return defaultSelectorFactory(store.dispatch, selectorFactoryOptions)
      }, [store])

      const [subscription, notifyNestedSubs] = React.useMemo(() => {
        const subscription = createSubscription(store, contextValue.subscription)

        const notifyNestedSubs = subscription.notifyNestedSubs.bind(subscription)

        return [subscription, notifyNestedSubs]
      }, [store, contextValue])

      const overriddenContextValue = React.useMemo(() => {
        return {
          ...contextValue,
          subscription,
        }
      }, [contextValue, subscription])

      const lastChildProps = React.useRef()
      const lastProps = React.useRef()
      const childPropsFromStoreUpdate = React.useRef()
      const renderIsScheduled = React.useRef(false)
      const isMounted = React.useRef(false)

      const latestSubscriptionCallbackError = React.useRef()

      useIsomorphicLayoutEffect(() => {
        isMounted.current = true
        return () => {
          isMounted.current = false
        }
      }, [])

      const actualChildPropsSelector = React.useMemo(() => {
        const selector = () => {
          if (
            childPropsFromStoreUpdate.current &&
            props === lastProps.current
          ) {
            return childPropsFromStoreUpdate.current
          }

          // pass state to mapStateToProps
          return childPropsSelector(store.getState(), props)
        }

        return selector
      }, [store, props])

      const subscribeForReact = React.useMemo(() => {
        const subscribe = reactListener => {
          if (!subscription) {
            return () => {}
          }

          return subscribeUpdates(
            store,
            subscription,
            childPropsSelector,
            lastProps,
            lastChildProps,
            renderIsScheduled,
            isMounted,
            childPropsFromStoreUpdate,
            notifyNestedSubs,
            reactListener,
          )
        }

        return subscribe
      }, [subscription])

      useIsomorphicLayoutEffect(captureWrapperProps, [
        lastProps,
        lastChildProps,
        renderIsScheduled,
        props,
        childPropsFromStoreUpdate,
        notifyNestedSubs,
      ])

      // props after all mapToState wrappers ++ own props
      let actualChildProps = {}
      try {
        actualChildProps = React.useSyncExternalStore(
          subscribeForReact,
          actualChildPropsSelector,
        )
      } catch (err) {
        console.error(err)
        throw err
      }

      useIsomorphicLayoutEffect(() => {
        latestSubscriptionCallbackError.current = undefined
        childPropsFromStoreUpdate.current = undefined
        lastChildProps.current = actualChildProps
      })

      const renderedWrappedComponent = React.useMemo(() => {
        return (
          <WrappedComponent {...actualChildProps} />
        )
      }, [WrappedComponent, actualChildProps])


      const renderedChild = React.useMemo(() => {
        return (
          <ReactReduxContext.Provider value={overriddenContextValue}>
            {renderedWrappedComponent}
          </ReactReduxContext.Provider>
        )
      }, [ReactReduxContext, renderedWrappedComponent, overriddenContextValue])

      return renderedChild
    }

    const Connect = React.memo(ConnectFunction)
    Connect.WrappedComponent = WrappedComponent
    Connect.displayName = ConnectFunction.displayName = displayName

    /* TODO: forwardRef */

    return hoistStatics(Connect, WrappedComponent)
  }

  return wrapWithConnect;
}