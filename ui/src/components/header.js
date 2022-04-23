import { NavLink } from 'react-router-dom'
import styles from './header.module.css'
import logo64 from '../images/logo64.png'

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
            <a href="#" >Logout</a>
          </div>
        </div>
      </div>
    </div>
  )
}
