import { Outlet } from 'react-router-dom'
import Header from './components/header'
import { WebSocketProvider } from './lib/websocket'

export default function App() {
  return (
    <WebSocketProvider>
      <Header />
      <div className="container">
        <Outlet />
      </div>
    </WebSocketProvider>
  )
}
