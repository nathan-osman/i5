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
      render: row => row.domain
    },
    {
      title: "Status",
      render: row => row.running ?
        <div className={styles.running}>Running</div> :
        <div className={styles.stopped}>Stopped</div>
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
