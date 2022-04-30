import { useEffect, useMemo } from 'react'
import { usePopup } from '../lib/popup'

export default function Requests() {

  const popup = usePopup()

  const webSocket = useMemo(() => {
    const secure = location.protocol.startsWith('https')
    const webSocket = new WebSocket(
      `${secure ? 'wss' : 'ws'}://${location.host}/api/ws`
    )
    webSocket.onopen = () => {
      //...
    }
    webSocket.onerror = (e) => {
      popup.error(e.message)
    }
    webSocket.onclose = () => {
      popup.info("WebSocket connection lost.")
    }
    return webSocket
  }, [])

  useEffect(() => {
    return () => {
      webSocket.close()
    }
  }, [])

  return (
    <>
      <div className="title">Requests</div>
    </>
  )
}
