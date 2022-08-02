import ReactTimeAgo from 'react-time-ago'
import Table from '../components/table'
import { useApi } from '../lib/api'

export default function Home() {

  const api = useApi()

  const headers = [
    {
      title: "Name",
      render: row => row.name
    },
    {
      title: "Value",
      render: row => row.value
    }
  ]

  const rows = [
    {
      name: "Go Version",
      value: api.status.go_version
    },
    {
      name: "Startup",
      value: <ReactTimeAgo date={api.status.startup * 1000} />
    }
  ]

  return (
    <>
      <div className="title">Home</div>
      <div className="secondary">
        i5 provides this web interface for interacting with containers and services from the host. You can also use this interface for monitoring services and requests, watching traffic in realtime.
      </div>
      <div className="subtitle">Status</div>
      <div className="secondary">
        Information about the host and i5 server are shown below:
      </div>
      <div className="half">
        <Table
          headers={headers}
          rows={rows}
        />
      </div>
    </>
  )
}
