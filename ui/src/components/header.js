import styles from './header.module.css'
import logo64 from '../images/logo64.png'

export default function Header() {
  return (
    <div className={styles.header}>
      <div className="container">
        <div className={styles.links}>
          <img src={logo64} />
          <a href="/">Containers</a>
          <a href="/">Requests</a>
        </div>
      </div>
    </div>
  )
}
