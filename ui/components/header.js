import styles from './header.module.css'
import { useRouter } from 'next/router'

function ActiveLink({ children, href }) {

  const router = useRouter()
  const isActive = href === router.asPath

  const handleClick = (e) => {
    e.preventDefault()
    router.push(href)
  }

  return (
    <a
      href={href}
      onClick={handleClick}
      className={isActive ? styles.active : ''}>
      {children}
    </a>
  )
}

export default function Header() {
  return (
    <div className={styles.header}>
      <div className="container">
        <div className={styles.links}>
          <img
            src="/logo192.png"
            width="32px" />
          <ActiveLink href="/containers">Containers</ActiveLink>
          <ActiveLink href="/requests">Requests</ActiveLink>
          <div className={styles.spacer} />
          <ActiveLink href="/logout">Logout</ActiveLink>
        </div>
      </div>
    </div>
  )
}
