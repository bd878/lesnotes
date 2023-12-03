import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';

const LoginForm = lazy(() => import("../../components/LoginForm/index.jsx"));

const Login = () => (
  <Suspense fallback={i18n('loading')}>
    <div>{i18n("login_form_header")}</div>

    <LoginForm />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
