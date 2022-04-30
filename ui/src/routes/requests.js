import RequestList from '../components/requestlist'
import { WebSocketProvider } from '../lib/websocket'

export default function Requests() {
  return (
    <WebSocketProvider>
      <div className="title">Requests</div>
      <div className="secondary">
        This page enables you to watch incoming requests to the i5 server in realtime. Detailed information about the inidividual requests are included as well as aggregate information that accumulates while this page is open.
      </div>
      <RequestList />
    </WebSocketProvider>
  )
}
