import {
  createContext,
  useContext,
  useEffect,
  useMemo
} from 'react'

const WebSocketContext = createContext(null)

function WebSocketProvider({ children }) {

  function connectToWebSocket() {
    const secure = location.protocol.startsWith('https')
    return new WebSocket(
      `${secure ? 'wss' : 'ws'}://${location.host}/api/ws`
    )
  }

  // Create the websocket when mounting
  const webSocket = useMemo(() => connectToWebSocket(), [])

  // Close the websocket when unmounting
  useEffect(() => {
    return () => {
      webSocket.close()
    }
  })

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
