import React from 'react';
import ReactDOM from 'react-dom/client';
import i18n from '../../i18n';

const XXX = () => (
  <div>{i18n("not_found")}</div>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<XXX />);
