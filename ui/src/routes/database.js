import { useParams } from 'react-router-dom'
import { useApi } from '../lib/api'
import DatabaseList from '../components/databaselist'

export default function Database() {

  const api = useApi()

  let { name } = useParams()

  let database = api.status.databases[name]
  if (database === undefined) {
    return (
      <>
        <div className="title">Not Found</div>
        <div className="secondary">Database "{name}" does not exist.</div>
      </>
    )
  }

  return (
    <>
      <div className="title">{database.title}</div>
      <div className="secondary">
        <strong>Version:</strong>{' '}
        {database.version}
      </div>
      <DatabaseList database={database} />
    </>
  )
}
