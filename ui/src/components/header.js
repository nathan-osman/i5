import { NavLink, useNavigate } from 'react-router-dom'
import { useApi } from '../lib/api'
import logo64 from '../images/logo64.png'
import styles from './header.module.css'

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
        // TODO: use new message interface
      })
  }

  return (
    <div className={styles.header_outer}>
      <div className="container">
        <div className={styles.header_inner}>
          <NavLink to="/">
            <img src={logo64} className={styles.logo} />
          </NavLink>
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
