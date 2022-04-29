import { useEffect, useState } from 'react'
import { useApi } from '../lib/api'
import Spinner from './spinner'
import styles from './data.module.css'

export default function Data({ children, url, onData }) {

  const api = useApi()

  const [isLoading, setIsLoading] = useState(true)
  const [errorMessage, setErrorMessage] = useState(null)

  useEffect(() => {
    api.fetch(url)
      .then(onData)
      .catch((e) => {
        setErrorMessage(e.message)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }, [])

  if (isLoading) {
    return (
      <div className={styles.block}>
        <Spinner />
      </div>
    )
  }

  if (errorMessage !== null) {
    return (
      <div className={styles.block}>
        <span className={styles.error}>
          Error: {errorMessage}
        </span>
      </div>
    )
  }

  return children
}
