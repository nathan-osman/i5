import { useEffect, useState } from 'react'
import prettyBytes from 'pretty-bytes'
import { useWebSocket } from '../lib/websocket'
import Client from './client'
import Table from './table'
import styles from './requestlist.module.css'
import Stat from './stat'

export default function RequestList() {

  const headers = [
    {
      title: "Time",
      render: row => <>
        {row.time.getHours().toString().padStart(2, '0')}:
        {row.time.getMinutes().toString().padStart(2, '0')}:
        {row.time.getSeconds().toString().padStart(2, '0')}
        <small className="secondary">
          .{row.time.getMilliseconds().toString().padStart(3, '0')}
        </small>
      </>
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
      expand: true,
      render: row => {
        return (
          <div className={styles.request}>
            <div className={styles.method}>{row.method}</div>
            <div className={styles.host}>{row.host}</div>
            <div className={styles.path} title={row.path}>{row.path}</div>
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
        } else if (row.status_code >= 400) {
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
      right: true,
      render: row => prettyBytes(parseInt(row.content_length) || 0)
    }
  ]

  const webSocket = useWebSocket()

  const [numRequests, setNumRequests] = useState(0)
  const [bandwidth, setBandwidth] = useState(0)
  const [requests, setRequests] = useState([])

  useEffect(() => {
    function processResponse(e) {
      setNumRequests(numRequests => numRequests + 1)
      try {
        setBandwidth(bandwidth => {
          return bandwidth + (parseInt(e.detail.content_length) || 0)
        })
      } catch { }
      setRequests((requests) => [
        {
          ...e.detail,
          time: new Date()
        },
        ...requests
      ].slice(0, 16))
    }
    webSocket.addEventListener('response', processResponse)
    return () => {
      webSocket.removeEventListener('response', processResponse)
    }
  }, [])

  return (
    <>
      <div className={styles.stats}>
        <Stat
          title="request(s)"
          value={numRequests}
        />
        <Stat
          title="bandwidth"
          value={prettyBytes(bandwidth)}
        />
      </div>
      <Table
        headers={headers}
        rows={requests}
      />
    </>
  )
}
