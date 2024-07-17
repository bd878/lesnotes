import React from 'react';
import ReactDOM from 'react-dom/client';
import Tag from '../../components/Tag';
import i18n from '../../i18n';

const XXX = () => (
  <Tag>{i18n("not_found")}</Tag>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<XXX />);
