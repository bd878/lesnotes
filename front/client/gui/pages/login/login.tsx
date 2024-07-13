import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';
import "./index.sass";

const LoginForm = lazy(() => import("../../components/LoginForm"));

const Login = () => (
  <Suspense fallback={i18n('loading')}>
    <div className="login_page">{i18n("login_form_header")}</div>

    <LoginForm />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
