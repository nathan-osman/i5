import { useState } from 'react'
import Data from './data'
import Table from './table'
import Tool from './tool'
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
          {row.running ?
            <span className={styles.running}>Running</span> :
            <span className={styles.stopped}>Stopped</span>
          }
        </div>
      )
    },
    {
      title: "Tools",
      right: true,
      render: row => row.running ?
        <Tool
          src={stopIcon}
          url={`/api/containers/${row.id}/state`}
          data={{ action: 'stop' }}
        /> :
        <Tool
          src={startIcon}
          url={`/api/containers/${row.id}/state`}
          data={{ action: 'start' }}
        />
    }
  ]

  const [containers, setContainers] = useState([])

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
