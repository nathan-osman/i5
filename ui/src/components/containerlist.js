import { useEffect, useState } from 'react'
import ReactTimeAgo from 'react-time-ago'
import Data from './data'
import Table from './table'
import Tool from './tool'
import { useWebSocket } from '../lib/websocket'
import styles from './containerlist.module.css'
import startIcon from '../images/icons/start.svg'
import stopIcon from '../images/icons/stop.svg'

export default function ContainerList() {

  const headers = [
    {
      title: "Name",
      render: row => row.title
    },
    {
      title: "Domains",
      expand: true,
      render: row => row.domains.map(domain =>
        <div key={domain}>
          <a href={`https://${domain}`} target="_blank">
            {domain}
          </a>
        </div>
      )
    },
    {
      title: "Uptime",
      render: row => row.uptime ?
        <ReactTimeAgo date={row.uptime * 1000} /> :
        <></>
    },
    {
      title: "Status",
      render: row => (
        <div className={styles.status}>
          {row.running ?
            <span className={styles.running}>Running</span> :
            <span className={styles.stopped}>Stopped</span>
          }
        </div>
      )
    },
    {
      title: "Tools",
      render: row => (
        <>
          {
            row.running ?
              <Tool
                src={stopIcon}
                url={`/api/containers/${row.id}/state`}
                data={{ action: 'stop' }}
                title="Stop the container"
              /> :
              <Tool
                src={startIcon}
                url={`/api/containers/${row.id}/state`}
                data={{ action: 'start' }}
                title="Start the container"
              />
          }
        </>
      )
    }
  ]

  const webSocket = useWebSocket()

  const [containers, setContainers] = useState([])

  useEffect(() => {
    function processContainerAction(e) {
      const container = e.detail.container
      switch (e.detail.action) {
        case 'create':
          setContainers(containers => [...containers, container])
          break
        case 'destroy':
          setContainers(containers => containers.filter(c => c.id != container.id))
          break
        case 'start':
        case 'die':
          setContainers(containers => containers.map(c => {
            return c.id == container.id ? container : c
          }))
          break
      }
    }
    webSocket.addEventListener('container', processContainerAction)
    return () => {
      webSocket.removeEventListener('container', processContainerAction)
    }
  }, [])

  function handleData(d) {
    d.sort((a, b) => {
      const titleA = a.title.toUpperCase()
      const titleB = b.title.toUpperCase()
      if (titleA < titleB) { return -1 }
      if (titleB < titleA) { return 1 }
      return 0
    })
    setContainers(d)
  }

  return (
    <Data url="/api/containers" onData={handleData}>
      <Table
        headers={headers}
        rows={containers}
      />
    </Data>
  )
}
