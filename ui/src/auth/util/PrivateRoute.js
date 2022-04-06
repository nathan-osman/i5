import { Navigate } from 'react-router-dom';
import { useRecoilValue } from 'recoil';
import { authAtom } from '../api/auth';

/**
 * Component to ensure user is logged in before navigating to a page
 */
const PrivateRoute = ({ children }) => {

  const auth = useRecoilValue(authAtom);

  // If the user is authenticated, render the page content; otherwise, redirect
  // the user to the login page
  return auth.isAuthenticated ? children : <Navigate to="/login" />;
};

export default PrivateRoute;
