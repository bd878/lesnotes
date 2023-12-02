import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';

const RegisterForm = lazy(() => import("../../components/RegisterForm/index.jsx"));

const Register = () => (
  <Suspense fallback="Loading...">
    <div>Register:</div>

    <RegisterForm />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Register />);
