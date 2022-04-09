import axios from 'axios';
import 'bootstrap';
import React from 'react';
import ReactDOM from 'react-dom';
import { RecoilRoot } from 'recoil';
import { HashRouter, Routes, Route } from "react-router-dom";

import App from './app/App';
import Home from './pages/Home';
import Login from './auth/Login';
import PrivateRoute from './auth/util/PrivateRoute';
import ContainerList from './components/ContainerList';
import RequestTicker from './components/RequestTicker';
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.css';

if (typeof process.env.REACT_APP_HOSTNAME !== 'undefined') {
  axios.defaults.baseURL = `http://${process.env.REACT_APP_HOSTNAME}`;
  axios.defaults.withCredentials = true;
}

ReactDOM.render(
  <React.StrictMode>
    <RecoilRoot>
      <HashRouter>
        <Routes>
          <Route path="/" element={
            <PrivateRoute>
              <App />
            </PrivateRoute>
          }>
            <Route index element={<Home />} />
            <Route path="/containers" element={<ContainerList />} />
            <Route path="/requests" element={<RequestTicker />} />
          </Route>
          <Route path="/login" element={<Login />} />
        </Routes>
      </HashRouter>
    </RecoilRoot>
  </React.StrictMode>,

  document.getElementById('root')
);
