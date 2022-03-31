import axios from 'axios';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSetRecoilState } from 'recoil';
import { authState } from '../api/state';
import './index.css';

const Login = () => {

  let navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const setAuthState = useSetRecoilState(authState);

  const login = async (username, password) => {
    try {
      await axios.post('/auth/login', { username, password });
      setAuthState((authState) => {
        return { ...authState, isAuthenticated: true };
      });
      navigate("/");
    } catch (e) {
      //...
    } finally {
      //...
    }
  };

  const onSubmit = (e) => {
    e.preventDefault();
    login(username, password);
  };

  const onUsernameChange = ({ target: { value } }) => {
    setUsername(value);
  };

  const onPasswordChange = ({ target: { value } }) => {
    setPassword(value);
  };

  return (
    <div id="Login" className="bg-secondary bg-gradient">
      <div className="card text-dark bg-light mb-3" style={{ width: '18rem' }}>
        <div className="card-header">Login</div>
        <div className="card-body">
          <form onSubmit={onSubmit}>
            <div className="mb-3">
              <label htmlFor="username" className="form-label">Username</label>
              <input
                type="text"
                className="form-control"
                id="username"
                value={username}
                onChange={onUsernameChange}
                autoFocus />
            </div>
            <div className="mb-3">
              <label htmlFor="password" className="form-label">Password</label>
              <input
                type="password"
                className="form-control"
                id="password"
                value={password}
                onChange={onPasswordChange} />
            </div>
            <button type="submit" className="btn btn-primary">Submit</button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
