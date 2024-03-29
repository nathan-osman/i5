import { useEffect, useState } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import Header from './header'
import Splash from './splash'
import Spinner from './spinner'
import { useApi } from '../lib/api'
import { PopupProvider } from '../lib/popup'
import { WebSocketProvider } from '../lib/websocket'

export default function Page() {

  const api = useApi()
  const navigate = useNavigate()

  const [isLoggingIn, setIsLoggingIn] = useState(true)

  useEffect(() => {
    api.fetchStatus()
      .then(() => {
        setIsLoggingIn(false)
      })
      .catch(() => {
        navigate(`/login?url=${location.pathname}`)
      })
  }, [])

  if (isLoggingIn) {
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
