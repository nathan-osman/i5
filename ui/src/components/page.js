import { useApi } from '../lib/api'
import { Navigate } from 'react-router-dom'
import Header from './header'

export default function Page({ children }) {

  const api = useApi()

  if (!api.isLoggedIn) {
    return <Navigate to="/login" />
  }

  return (
    <>
      <Header />
      <div className="container">
        {children}
      </div>
    </>
  )
}
