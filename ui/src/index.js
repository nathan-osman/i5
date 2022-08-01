import '@fontsource/source-sans-pro'
import 'core-js/stable'
import 'regenerator-runtime/runtime'
import { createRoot } from 'react-dom/client'
import {
  BrowserRouter,
  Routes,
  Route
} from 'react-router-dom'
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en.json'
import App from './app'
import Page from './components/page'
import Splash from './components/splash'
import Home from './routes/home'
import Containers from './routes/containers'
import Requests from './routes/requests'
import Database from './routes/database'
import Login from './routes/login'
import './index.css'

TimeAgo.addDefaultLocale(en)

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
          <Route path="/db/:name" element={<Database />} />
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
