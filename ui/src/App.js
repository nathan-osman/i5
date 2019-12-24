import React from 'react';
import './App.css';

import ContainerTable from './components/ContainerTable';

function App() {
  return (
    <div className="App">

      {/* Application header */}
      <div className="App-header">
        <div className="container">
          i5 Status
        </div>
      </div>

      <div className="container">
        <ContainerTable containers={[]} />
      </div>

    </div>
  );
}

export default App;
