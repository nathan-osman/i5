import '@fontsource/source-sans-pro'
import '@fontsource/source-sans-pro/300.css'
import 'core-js/stable'
import 'regenerator-runtime/runtime'
import { createRoot } from 'react-dom/client'
import {
  BrowserRouter,
  Routes,
  Route
} from 'react-router-dom'
import App from './app'
import Page from './components/page'
import Splash from './components/splash'
import Home from './routes/home'
import Containers from './routes/containers'
import Requests from './routes/requests'
import Login from './routes/login'
import './index.css'

const root = createRoot(
  document.getElementById('root')
)

root.render(
  <BrowserRouter>
    <Routes>
      <Route element={<App />}>
        <Route element={<Page />}>
          <Route path="/" element={<Home />} />
          <Route path="/containers" element={<Containers />} />
          <Route path="/requests" element={<Requests />} />
        </Route>
        <Route path="/login" element={<Login />} />
        <Route path="/*" element={
          <Splash>
            <div className="title">Not Found</div>
            <div className="secondary">The page you tried to view does not exist.</div>
          </Splash>
        } />
      </Route>
    </Routes>
  </BrowserRouter>
)
