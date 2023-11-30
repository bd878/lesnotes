import React, {Suspense, lazy} from 'react';
import ReactDOM from 'react-dom/client';

const List = lazy(() => import("../../components/List/index.jsx"));
const ListItem = lazy(() => import("../../components/ListItem/index.jsx"));

const Messages = () => (
  <Suspense fallback="Loading...">
    <div>Messages:</div>

    <List>
      <ListItem>Message 1</ListItem>
      <ListItem>Message 2</ListItem>
    </List>
  </Suspense>
)

const root = ReactDOM.createRoot(document.getElementById('app'));
root.render(<Messages />);
