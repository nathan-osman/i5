import axios from 'axios';
import 'bootstrap';
import React from 'react';
import ReactDOM from 'react-dom';
import { RecoilRoot } from 'recoil';
import { HashRouter, Routes, Route } from "react-router-dom";

import App from './app/App';
import Containers from './pages/Containers';
import Home from './pages/Home';
import Login from './auth/Login';
import PrivateRoute from './auth/util/PrivateRoute';
import Requests from './pages/Requests';

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
            <Route path="/containers" element={<Containers />} />
            <Route path="/requests" element={<Requests />} />
          </Route>
          <Route path="/login" element={<Login />} />
        </Routes>
      </HashRouter>
    </RecoilRoot>
  </React.StrictMode>,

  document.getElementById('root')
);
