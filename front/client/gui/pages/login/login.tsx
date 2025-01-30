import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import StoreProvider from '../../providers/Store';
import i18n from '../../i18n';

const LoginForm = lazy(() => import("../../components/LoginForm"));

const Login = () => (
  <Suspense fallback={i18n('loading')}>
    <Tag css="flex white bg-primary">
      {i18n("login_form_header")}
    </Tag>

    <StoreProvider>
      <LoginForm />
    </StoreProvider>
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
