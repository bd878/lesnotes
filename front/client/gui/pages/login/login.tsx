import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import StoreProvider from '../../providers/Store';
import AuthProvider from '../../providers/Auth';
import i18n from '../../i18n';

const LoginForm = lazy(() => import("../../components/LoginForm"));

function Login() {
  return (
    <Suspense fallback={i18n('loading')}>
      <StoreProvider>
        <AuthProvider inverted={true}>
          <Tag css="flex white bg-primary">
            {i18n("login_form_header")}
          </Tag>

          <LoginForm />
        </AuthProvider>
      </StoreProvider>
    </Suspense>
  )
}

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Login />);
