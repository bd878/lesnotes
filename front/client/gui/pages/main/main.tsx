import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';

const Main = () => (
  <>
    <Tag>{i18n("index_intro")}</Tag>
    <Tag el="a" href="/login">{i18n("login")}</Tag>
    <Tag el="a" href="/register">{i18n("register")}</Tag>
  </>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Main />);
