import { useState } from 'react'
import Table from './table'

export default function RequestList() {

  const headers = [
    {
      title: "Time"
    },
    {
      title: "Method"
    },
    {
      title: "Path"
    },
    {
      title: "Type"
    },
    {
      title: "Size"
    }
  ]

  const [requests, setRequests] = useState([])

  return (
    <>
      <Table
        headers={headers}
        rows={requests}
      />
    </>
  )
}
