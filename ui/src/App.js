import React from 'react';
import './App.css';

import ContainerList from './components/ContainerList';

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
        <ContainerList />
      </div>

    </div>
  );
}

export default App;
