import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';

const Form = lazy(() => import("../../components/Form/index.jsx"));
const FormField = lazy(() => import("../../components/FormField/index.jsx"));

const Register = () => (
  <Suspense fallback="Loading...">
    <div>Register:</div>

    <Form>
      <FormField name="name" type="text" />
      <FormField name="password" type="password" />
    </Form>
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Register />);
