import { useState } from 'react'
import { useApi } from '../lib/api'
import { usePopup } from '../lib/popup'
import styles from './tool.module.css'
import spinner from '../images/spinner.svg'

export default function Tool({ src, url, data }) {

  const api = useApi()
  const popup = usePopup()

  const [isLoading, setIsLoading] = useState(false)

  function handleClick() {
    setIsLoading(true)
    api.fetch(url, data)
      .catch(e => popup.error(e.message))
      .finally(() => setIsLoading(false))
  }

  return (
    <img className={styles.tool} src={isLoading ? spinner : src} onClick={handleClick} />
  )
}
