import { Outlet } from 'react-router-dom'
import { ApiProvider } from './lib/api'

export default function App() {
  return (
    <ApiProvider>
      <Outlet />
    </ApiProvider>
  )
}
