import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';

const LoginForm = lazy(() => import("../../components/LoginForm/index.jsx"));

const Login = () => (
  <Suspense fallback="Loading...">
    <div>Login:</div>

    <LoginForm />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
