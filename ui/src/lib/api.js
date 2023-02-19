import { createContext, useContext, useState } from 'react'

class HttpError extends Error {
  constructor(response, message) {
    super(message)
    this.response = response
  }
}

const ApiContext = createContext(null)

function ApiProvider({ children }) {

  const [isActive, setIsActive] = useState(false)
  const [status, setStatus] = useState(null)

  async function fetchInternal(url, method, data) {
    let init = { method }
    if (data != undefined) {
      init['headers'] = { 'Content-Type': 'application/json' }
      init['body'] = JSON.stringify(data)
    }
    const response = await fetch(url, init)
    if (!response.ok) {
      let errorMessage
      try {
        errorMessage = (await response.json()).error
      } catch (e) {
        throw new HttpError(response, `${response.status} ${response.statusText}`)
      }
      throw new HttpError(response, errorMessage)
    }
    return await response.json()
  }

  const api = {
    isActive,
    status,
    fetch: async (url, method, data) => {
      setIsActive(true)
      try {
        return await fetchInternal(url, method, data)
      } finally {
        setIsActive(false)
      }
    },
    fetchStatus: async () => {
      setStatus(await api.fetch('/api/status'))
    },
  }

  return (
    <ApiContext.Provider value={api}>
      {children}
    </ApiContext.Provider>
  )
}

function useApi() {
  return useContext(ApiContext)
}

export { ApiProvider, useApi }
