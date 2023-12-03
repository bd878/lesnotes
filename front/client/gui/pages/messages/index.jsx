import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import AuthProvider from '../../providers/Auth';
import i18n from '../../i18n';

const MessagesList = lazy(() => import("../../components/MessagesList/index.jsx"));
const SendMessageForm = lazy(() => import("../../components/SendMessageForm/index.jsx"));

const Messages = () => (
  <AuthProvider fallback={i18n("messages_auth_fallback")}>
    <Suspense fallback={i18n("loading")}>
      <div>
        <MessagesList />

        <SendMessageForm />
      </div>
    </Suspense>
  </AuthProvider>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
