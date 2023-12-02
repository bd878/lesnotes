import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import AuthProvider from '../../providers/Auth';

const MessagesList = lazy(() => import("../../components/MessagesList/index.jsx"));
const SendMessageForm = lazy(() => import("../../components/SendMessageForm/index.jsx"));

const Messages = () => (
  <AuthProvider fallback="Go to /login">
    <Suspense fallback="Loading...">
      <div>
        <MessagesList />

        <SendMessageForm />
      </div>
    </Suspense>
  </AuthProvider>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
