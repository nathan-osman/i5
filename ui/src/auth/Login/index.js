import { useNavigate } from 'react-router-dom';
import { useSetRecoilState } from 'recoil';
import { authAtom } from '../api/auth';
import Form from '../../components/Form';
import './index.css';

const Login = () => {

  const navigate = useNavigate();
  const setAuth = useSetRecoilState(authAtom);

  function handleData() {
    setAuth((auth) => {
      return { ...auth, isAuthenticated: true };
    });
    navigate("/");
  }

  return (
    <div id="Login" className="bg-secondary bg-gradient">
      <div className="card text-dark bg-light mb-3" style={{ width: '18rem' }}>
        <div className="card-header">Login</div>
        <div className="card-body">
          <Form url="/auth/login" onData={handleData} id="Login">
            <div className="mb-3">
              <label htmlFor="username" className="form-label">Username</label>
              <input
                type="text"
                className="form-control"
                id="username"
                name="username"
                autoFocus />
            </div>
            <div className="mb-3">
              <label htmlFor="password" className="form-label">Password</label>
              <input
                type="password"
                className="form-control"
                id="password"
                name="password" />
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
};

export default Login;
