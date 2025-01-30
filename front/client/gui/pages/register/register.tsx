import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../i18n';
import StoreProvider from '../../providers/Store';

const RegisterForm = lazy(() => import("../../components/RegisterForm"));

const Register = () => (
  <Suspense fallback={i18n('loading')}>
    <Tag>{i18n("register")}</Tag>

    <StoreProvider>
      <RegisterForm />
    </StoreProvider>
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Register />);
