import spinner from '../images/spinner.svg'
import styles from './spinner.module.css'

export default function Spinner() {
  return (
    <div className={styles.spinner}>
      <img src={spinner} />
      <div>
        <div className={styles.loading}>Loading data</div>
        <div>Please wait...</div>
      </div>
    </div>
  )
}
