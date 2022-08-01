import { Outlet } from 'react-router-dom'
import Header from './header'
import Splash from './splash'
import Spinner from './spinner'
import { useApi } from '../lib/api'
import { PopupProvider } from '../lib/popup'
import { WebSocketProvider } from '../lib/websocket'

export default function Page() {

  const api = useApi()

  if (api.isLoggingIn) {
    return (
      <Splash><Spinner /></Splash>
    )
  }

  return (
    <PopupProvider>
      <WebSocketProvider>
        <Header />
        <div className="container">
          <Outlet />
        </div>
      </WebSocketProvider>
    </PopupProvider>
  )
}
