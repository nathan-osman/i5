import { useState } from 'react'
import Data from './data'
import Table from './table'

export default function DatabaseList({ database }) {

  const headers = [
    {
      title: "Database Name",
      render: row => row
    }
  ]

  const [databases, setDatabases] = useState([])

  function handleData(d) {
    d.sort()
    setDatabases(d)
  }

  return (
    <>
      <div className="subtitle">Databases</div>
      <div className="secondary">
        This table displays a list of all {database.title} databases.
      </div>
      <Data
        url={`/api/db/${database.name}/databases`}
        onData={handleData}
        dependencies={[database]}
      >
        <Table
          headers={headers}
          rows={databases}
        />
      </Data>
    </>
  )
}
