import styles from './splash.module.css'

export default function Splash({ children }) {
  return (
    <div className={styles.splash}>
      <div>{children}</div>
    </div>
  )
}
