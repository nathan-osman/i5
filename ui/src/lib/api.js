import { createContext, useContext, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

class HttpError extends Error {
  constructor(response, message) {
    super(message)
    this.response = response
  }
}

const ApiContext = createContext(null)

function ApiProvider({ children }) {

  const navigate = useNavigate()

  const [isActive, setIsActive] = useState(false)
  const [isLoggingIn, setIsLoggingIn] = useState(true)
  const [status, setStatus] = useState(null)

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
    isLoggingIn,
    status,
    fetch: async (url, data) => {
      setIsActive(true)
      try {
        return await fetchInternal(url, data)
      } finally {
        setIsActive(false)
      }
    },
  }

  useEffect(() => {
    api.fetch('/api/status')
      .then(d => {
        setIsLoggingIn(false)
        setStatus(d)
      })
      .catch(() => {
        navigate(`/login?url=${location.pathname}`)
      })
  }, [])

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
