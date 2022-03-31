import { Link } from 'react-router-dom';

const Header = () => {
  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark mb-4">
      <div className="container">
        <Link className="navbar-brand" to="/">i5 Status</Link>
        <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbar-nav">
          <span className="navbar-toggler-icon"></span>
        </button>
        <div className="collapse navbar-collapse" id="navbar-nav">
          <div className="navbar-nav me-auto">
          </div>
          <div className="navbar-nav">
            <Link className="nav-link" to="/login">Login</Link>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Header;
