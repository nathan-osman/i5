import { useParams } from "react-router-dom"
import { useApi } from "../lib/api"


export default function Database() {

  const api = useApi()

  let { name } = useParams()

  let dbInfo = api.status.databases[name]

  return (
    <>
      <div className="title">{dbInfo.title}</div>
      <div className="secondary">{dbInfo.version}</div>
    </>
  )
}
