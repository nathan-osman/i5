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

  async function fetchInternal(url, data) {
    let init
    if (typeof data !== 'undefined') {
      init = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      }
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
    fetch: async (url, data) => {
      setIsActive(true)
      try {
        return await fetchInternal(url, data)
      } finally {
        setIsActive(false)
      }
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
