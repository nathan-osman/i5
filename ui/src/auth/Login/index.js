import './index.css';

const Login = () => {
  return (
    <div id="Login" className="bg-secondary bg-gradient">
      <div class="card text-dark bg-light mb-3" style={{ width: '18rem' }}>
        <div class="card-header">Login</div>
        <div class="card-body">
          <form>
            <div className="mb-3">
              <label for="username" className="form-label">Username</label>
              <input type="text" className="form-control" id="username" autoFocus />
            </div>
            <div className="mb-3">
              <label for="password" className="form-label">Password</label>
              <input type="password" className="form-control" id="password" />
            </div>
            <button type="submit" className="btn btn-primary">Submit</button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default Login;
