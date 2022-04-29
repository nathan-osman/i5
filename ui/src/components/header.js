import { NavLink, useNavigate } from 'react-router-dom'
import styles from './header.module.css'
import logo64 from '../images/logo64.png'
import { useApi } from '../lib/api'

function ActiveLink({ children, href }) {
  return (
    <NavLink
      to={href}
      className={({ isActive }) => isActive ? styles.active : null}
    >
      {children}
    </NavLink>
  )
}

export default function Header() {

  const api = useApi()
  const navigate = useNavigate()

  function handleLogout(e) {
    e.preventDefault()
    api.fetch('/auth/logout', {})
      .then(() => {
        navigate('/login')
      })
      .catch((e) => {
        //...
      })
  }

  return (
    <div className={styles.header_outer}>
      <div className="container">
        <div className={styles.header_inner}>
          <img src={logo64} className={styles.logo} />
          <div className={styles.nav}>
            <ActiveLink href="/containers">Containers</ActiveLink>
            <ActiveLink href="/requests">Requests</ActiveLink>
          </div>
          <div className={styles.separator} />
          <div className={styles.nav}>
            <a href="#" onClick={handleLogout}>Logout</a>
          </div>
        </div>
      </div>
    </div>
  )
}
