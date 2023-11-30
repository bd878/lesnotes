import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';

const Greeting = lazy(() => import("../../components/Greeting/index.jsx"));

const Index = () => (
  <Suspense fallback="Loading...">
    <Greeting />
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Index />);
