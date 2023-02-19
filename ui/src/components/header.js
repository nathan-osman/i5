import { useState } from 'react'
import { NavLink, useNavigate } from 'react-router-dom'
import { useApi } from '../lib/api'
import { usePopup } from '../lib/popup'
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
  const popup = usePopup()
  const navigate = useNavigate()

  const [isMenuOpen, setIsMenuOpen] = useState(false)

  function handleLogout(e) {
    e.preventDefault()
    api.fetch('/auth/logout', 'post')
      .then(() => {
        navigate('/login')
      })
      .catch((e) => {
        popup.error(e.message)
      })
  }

  function handleToggle() {
    setIsMenuOpen(isMenuOpen => !isMenuOpen)
  }

  function handleClose() {
    setIsMenuOpen(false)
  }

  return (
    <div className={styles.header_outer}>
      <div className="container">
        <div className={styles.header_inner}>
          <div className={styles.nav}>
            <NavLink to="/" onClick={handleClose}>
              <img src={logo64} className={styles.logo} />
            </NavLink>
            <div className={styles.toggle} onClick={handleToggle}>
              <div className={styles.i} />
              <div className={styles.i} />
              <div className={styles.i} />
            </div>
          </div>
          <div
            className={styles.menu_outer + (isMenuOpen ? '' : ` ${styles.hidden}`)}
            onClick={handleClose}
          >
            <div className={styles.menu_inner}>
              <ActiveLink href="/containers">Containers</ActiveLink>
              <ActiveLink href="/requests">Requests</ActiveLink>
              {
                Object.entries(api.status.databases).map(([k, v]) => (
                  <ActiveLink key={k} href={`/db/${k}`}>
                    {v.title}
                  </ActiveLink>
                ))
              }
            </div>
            <div className={styles.separator} />
            <div className={styles.menu_inner}>
              <a href="#" onClick={handleLogout}>
                Logout ({api.status.username})
              </a>
            </div>
          </div>
        </div>
      </div>
    </div >
  )
}
