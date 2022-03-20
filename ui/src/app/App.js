import { Outlet } from 'react-router-dom';

import Header from '../components/Header';

const App = () => {
  return (
    <div>

      {/* Header */}
      <Header />

      {/* Page content */}
      <div className="container mt-2">
        <Outlet />
      </div>

    </div >
  );
};

export default App;
