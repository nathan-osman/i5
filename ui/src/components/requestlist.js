import { useEffect, useState } from 'react'
import prettyBytes from 'pretty-bytes'
import { useWebSocket } from '../lib/websocket'
import Client from './client'
import Table from './table'
import styles from './requestlist.module.css'

export default function RequestList() {

  const headers = [
    {
      title: "Time",
      render: row => "-"
    },
    {
      title: "Client",
      render: row => <Client
        remoteAddr={row.remote_addr}
        countryCode={row.country_code}
        countryName={row.country_name}
      />
    },
    {
      title: "Request",
      render: row => {
        return (
          <div className={styles.request}>
            <div className={styles.method}>{row.method}</div>
            <div className={styles.host}>{row.host}</div>
            <div className={styles.path}>{row.path}</div>
          </div>
        )
      }
    },
    {
      title: "Response",
      render: row => {
        let className = styles.status_info
        if (row.status_code >= 200 && row.status_code < 300) {
          className = styles.status_good
        } else if (row.status_Code >= 400) {
          className = styles.status_bad
        }
        return <span className={`${styles.status} ${className}`}>{row.status}</span>
      }
    },
    {
      title: "Type",
      render: row => row.content_type || "-"
    },
    {
      title: "Size",
      render: row => {
        try {
          return prettyBytes(parseInt(row.content_length))
        } catch {
          return "-"
        }
      }
    }
  ]

  const webSocket = useWebSocket()

  const [requests, setRequests] = useState([])

  useEffect(() => {

    function processResponse(e) {
      setRequests((requests) => [e.detail, ...requests])
    }

    webSocket.addEventListener('response', processResponse)

    return () => {
      webSocket.removeEventListener('response', processResponse)
    }

  }, [])

  return (
    <>
      <Table
        headers={headers}
        rows={requests}
      />
    </>
  )
}
