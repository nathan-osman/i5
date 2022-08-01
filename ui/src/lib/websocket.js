import { createContext, useContext, useEffect, useMemo } from 'react'
import { usePopup } from './popup'

const WebSocketContext = createContext(null)

function WebSocketProvider({ children }) {

  const popup = usePopup()

  const eventTarget = useMemo(() => new EventTarget(), [])

  useEffect(() => {
    const secure = location.protocol.startsWith('https')

    let webSocket
    let reconnectTimerID

    function connect(reconnecting) {
      webSocket = new WebSocket(
        `${secure ? 'wss' : 'ws'}://${location.host}/api/ws`
      )
      webSocket.onopen = () => {
        if (reconnecting) {
          popup.info("reconnected to WebSocket")
        }
      }
      webSocket.onmessage = (e) => {
        const json = JSON.parse(e.data)
        const event = new CustomEvent(json.type, { detail: json.data })
        eventTarget.dispatchEvent(event)
      }
      webSocket.onerror = (e) => {
        popup.error(e.message)
      }
      webSocket.onclose = (e) => {
        if (!e.wasClean) {
          popup.error("WebSocket disconnected; reconnecting...")
          reconnectTimerID = setTimeout(() => connect(true), 30 * 1000)
        }
      }
    }

    // Make the initial connection
    connect(false)

    return () => {
      webSocket.close()
      clearTimeout(reconnectTimerID)
    }
  }, [])

  return (
    <WebSocketContext.Provider value={eventTarget}>
      {children}
    </WebSocketContext.Provider>
  )
}

function useWebSocket() {
  return useContext(WebSocketContext)
}

export { WebSocketProvider, useWebSocket }
