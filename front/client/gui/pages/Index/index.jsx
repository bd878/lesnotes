import React, {Suspense, lazy} from 'react';

const Greeting = lazy(() => import("../../components/Greeting/index.jsx"));

const Index = () => (
  <Suspense fallback="Loading...">
    <Greeting />
  </Suspense>
)

export default Index;
