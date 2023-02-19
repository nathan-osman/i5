import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import Form from '../components/form'
import Splash from '../components/splash'
import { useApi } from '../lib/api'

export default function Login() {

  const api = useApi()
  const location = useLocation()
  const navigate = useNavigate()

  const searchParams = new URLSearchParams(location.search)
  const url = searchParams.get('url') || '/'

  const [errorMessage, setErrorMessage] = useState(null)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  function handleSubmit(e) {
    e.preventDefault()
    api.fetch('/auth/login', 'post', { username, password })
      .then(() => {
        navigate(url)
      })
      .catch(e => {
        setErrorMessage(e.message)
      })
  }

  function handleUsernameChange(e) {
    setUsername(e.target.value)
  }

  function handlePasswordChange(e) {
    setPassword(e.target.value)
  }

  useEffect(() => {
    if (searchParams.has('url')) {
      setErrorMessage("You must login to access that page.")
    }
  }, [])

  return (
    <Splash>
      <Form
        title="Login"
        disabled={api.isActive}
        errorMessage={errorMessage}
        onSubmit={handleSubmit}
      >
        <div className="form-control">
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            onChange={handleUsernameChange}
            autoFocus />
        </div>
        <div className="form-control">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            onChange={handlePasswordChange} />
        </div>
      </Form>
    </Splash>
  )
}
