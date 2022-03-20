import { Container, Navbar } from 'react-bootstrap';

const Header = () => {
  return (
    <Navbar bg="dark" variant="dark">
      <Container>
        <Navbar.Brand>i5 Status</Navbar.Brand>
      </Container>
    </Navbar>
  );
};

export default Header;
