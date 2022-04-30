import { createContext, useContext, useEffect, useMemo } from 'react'
import { usePopup } from './popup'

const WebSocketContext = createContext(null)

function WebSocketProvider({ children }) {

  const popup = usePopup()

  const webSocket = useMemo(() => {
    const secure = location.protocol.startsWith('https')
    const webSocket = new WebSocket(
      `${secure ? 'wss' : 'ws'}://${location.host}/api/ws`
    )
    webSocket.onerror = (e) => {
      popup.error(e.message)
    }
    webSocket.onclose = (e) => {
      if (!e.wasClean) {
        popup.info("WebSocket connection lost")
      }
    }
    return webSocket
  }, [])

  useEffect(() => {
    return () => {
      webSocket.close()
    }
  }, [])

  return (
    <WebSocketContext.Provider value={webSocket}>
      {children}
    </WebSocketContext.Provider>
  )
}

function useWebSocket() {
  return useContext(WebSocketContext)
}

export { WebSocketProvider, useWebSocket }
