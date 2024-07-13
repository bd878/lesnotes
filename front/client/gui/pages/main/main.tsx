import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';

const Main = () => (
  <>
    <div>{i18n("index_intro")}</div>
    <a href="/login">{i18n("login")}</a>
    <a href="/register">{i18n("register")}</a>
  </>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
