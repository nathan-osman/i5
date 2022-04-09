import axios from 'axios';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

const Header = () => {

  const navigate = useNavigate();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  function handleLogout() {
    setIsLoggingOut(true);
    axios.post('/auth/logout')
      .then(() => navigate("/login"))
      .finally(() => setIsLoggingOut(false));
  }

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark mb-4">
      <div className="container">
        <Link className="navbar-brand" to="/">i5 Status</Link>
        <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbar-nav">
          <span className="navbar-toggler-icon"></span>
        </button>
        <div className="collapse navbar-collapse" id="navbar-nav">
          <div className="navbar-nav me-auto">
            <Link className="nav-link" to="/containers">Containers</Link>
            <Link className="nav-link" to="/requests">Requests</Link>
          </div>
          <div className="navbar-nav">
            {isLoggingOut ?
              <span className="navbar-text">Please wait...</span> :
              <button
                type="button"
                className="btn btn-dark"
                onClick={handleLogout}>
                Logout
              </button>
            }
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Header;
