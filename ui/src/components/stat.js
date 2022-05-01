import styles from './stat.module.css'

export default function Stat({ title, value }) {
  return (
    <div className={styles.stat}>
      <div className={styles.value}>{value}</div>
      <div className={styles.title}>{title}</div>
    </div>
  )
}
