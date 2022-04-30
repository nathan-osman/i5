import styles from './client.module.css'
import '../images/flags.css'

export default function Client({ remoteAddr, countryCode, countryName }) {
  return (
    <div className={styles.client}>
      <div className={`icon icon-${countryCode}`} />
      {remoteAddr}
    </div>
  )
}
