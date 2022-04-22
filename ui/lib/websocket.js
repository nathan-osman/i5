import { createContext, useContext, useEffect, useMemo } from 'react'

const WebSocketContext = createContext(null)

function WebSocketProvider({ children, url }) {

  function connectToWebSocket() {
    return new WebSocket(url)
  }

  const webSocket = useMemo(() => connectToWebSocket(), [])

  useEffect(() => {
    return () => {
      webSocket.close()
    }
  })

  return (
    <WebSocketContext.Provider value={webSocket}>{children}</WebSocketContext.Provider>
  )
}

function useWebSocket() {
  const context = useContext(WebSocketContext)
  return context
}

export { WebSocketProvider, useWebSocket }
