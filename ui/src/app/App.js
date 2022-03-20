import { Outlet } from 'react-router-dom'

import Header from '../components/Header';

const App = () => {
  return (
    <div>

      {/* Header */}
      <Header />

      {/* Page content */}
      <Outlet />

    </div>
  );
};

export default App;
