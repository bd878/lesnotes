import React, {useEffect, useState, Suspense} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';
import HomePage from '../../components/HomePage';
import Auth from '../../providers/Auth';
import {store} from '../../store';
import {
  addMessagesActionCreator,
  pushBackMessagesActionCreator,
  selectMessages,
} from '../../features/messages';

const StateContext = React.createContext(null)

function StoreProvider(props) {
  const [state, setState] = useState({})

  useEffect(() => {
    store.subscribe(function subscription() {
      try {
        setState(store.getState())
      } catch (e) {
        console.error(e)
      }
    })
  }, [setState])

  console.log(state)

  return (
    <StateContext.Provider value={{state}}>
      {props.children}
    </StateContext.Provider>
  )
}

const Home = () => {
  const dispatchAppendMessages = React.useCallback((messages) => {
    store.dispatch(addMessagesActionCreator(messages))
  }, [store.dispatch])

  const dispatchPushBackMessage = React.useCallback((messages) => {
    store.dispatch(pushBackMessagesActionCreator(messages))
  }, [store.dispatch])

  function ConnectedHomePage(props) {
    const contextValue = React.useContext(StateContext)

    return (
      <HomePage
        appendMessages={dispatchAppendMessages}
        pushBackMessages={dispatchPushBackMessage}
        messages={contextValue.state.messages ? selectMessages(contextValue.state) : []}
      />
    )
  }

  return (
    <Suspense fallback={i18n("loading")}>
      <Auth fallback={i18n("messages_auth_fallback")}>
        <StoreProvider>
          <ConnectedHomePage />
        </StoreProvider>
      </Auth>
    </Suspense>
  )
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Home />);
