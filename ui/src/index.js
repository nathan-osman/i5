import axios from 'axios';
import 'bootstrap';
import React from 'react';
import ReactDOM from 'react-dom';
import { RecoilRoot } from 'recoil';
import { BrowserRouter, Routes, Route } from "react-router-dom";

import App from './app/App';
import Home from './pages/Home';
import Login from './auth/Login';
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.css';

axios.defaults.baseURL = process.env.REACT_APP_API_PREFIX;

ReactDOM.render(
  <React.StrictMode>
    <RecoilRoot>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<App />}>
            <Route index element={<Home />} />
          </Route>
          <Route path="/login" element={<Login />} />
        </Routes>
      </BrowserRouter>
    </RecoilRoot>
  </React.StrictMode>,

  document.getElementById('root')
);
