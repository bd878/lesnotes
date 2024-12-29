import React from 'react';
import {Provider} from '../../third_party/react-redux'
import store from '../../store';

function StoreProvider(props) {
  return (
    <Provider store={store}>
      {props.children}
    </Provider>
  )
}

export default StoreProvider;