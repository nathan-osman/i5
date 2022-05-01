import { useState } from 'react'
import Data from './data'
import Table from './table'
import styles from './containerlist.module.css'

export default function ContainerList() {

  const headers = [
    {
      title: "Name",
      render: row => row.name
    },
    {
      title: "Domain",
      expand: true,
      render: row => (
        <a href={`https://${row.domain}`} target="_blank">
          {row.domain}
        </a>
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
      render: row => 'TODO'
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
