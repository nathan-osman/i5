import { useEffect, useState } from 'react'
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
      render: row => row.name
    },
    {
      title: "Domains",
      expand: true,
      render: row => row.domains.map(domain =>
        <div key={domain.id}>
          <a href={`https://${domain}`} target="_blank">
            {domain}
          </a>
        </div>
      )
    },
    {
      title: "Status",
      render: row => (
        <div className={styles.status}>
          {row.disabled ?
            <span className={styles.stopped}>Stopped</span> :
            <span className={styles.running}>Running</span>
          }
        </div>
      )
    },
    {
      title: "Tools",
      right: true,
      render: row => row.disabled ?
        <Tool
          src={startIcon}
          url={`/api/containers/${row.id}/state`}
          data={{ action: 'start' }}
        /> :
        <Tool
          src={stopIcon}
          url={`/api/containers/${row.id}/state`}
          data={{ action: 'stop' }}
        />
    }
  ]

  const webSocket = useWebSocket()

  const [containers, setContainers] = useState([])

  // TODO: sort containers by name

  useEffect(() => {
    function processContainerAction(e) {
      const container = e.detail.container
      switch (e.detail.action) {
        case 'create':
          setContainers(containers => [...containers, container])
          break;
        case 'destroy':
          setContainers(containers => containers.filter(c => c.id != container.id))
          break;
        case 'start':
        case 'die':
          setContainers(containers => containers.map(c => {
            return c.id == container.id ? container : c
          }))
          break;
      }
    }
    webSocket.addEventListener('container', processContainerAction)
    return () => {
      webSocket.removeEventListener('container', processContainerAction)
    }
  }, [])

  function handleData(d) {
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
