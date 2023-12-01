import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import AuthProvider from '../../providers/Auth';

const List = lazy(() => import("../../components/List/index.jsx"));
const ListItem = lazy(() => import("../../components/ListItem/index.jsx"));

const Messages = () => (
  <AuthProvider fallback="Authenticating...">
    <Suspense fallback="Loading...">
      <div>Messages:</div>

      <List>
        <ListItem>Message 1</ListItem>
        <ListItem>Message 2</ListItem>
      </List>
    </Suspense>
  </AuthProvider>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
