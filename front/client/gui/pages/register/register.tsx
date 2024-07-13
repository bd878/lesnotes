import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';

const RegisterForm = lazy(() => import("../../components/RegisterForm/RegisterForm.tsx"));

const Register = () => (
  <Suspense fallback={i18n('loading')}>
    <div>{i18n("register")}</div>

    <RegisterForm />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Register />);
