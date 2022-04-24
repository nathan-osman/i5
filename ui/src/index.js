import { createRoot } from 'react-dom/client'
import {
  BrowserRouter,
  Routes,
  Route
} from 'react-router-dom'
import App from './app'
import Containers from './routes/containers'
import Requests from './routes/requests'
import './index.css'

const root = createRoot(
  document.getElementById('root')
)

root.render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<App />}>
        <Route path="containers" element={<Containers />} />
        <Route path="requests" element={<Requests />} />
      </Route>
    </Routes>
  </BrowserRouter>,
)
