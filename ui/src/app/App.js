import { Container, Navbar } from 'react-bootstrap';
import ContainerList from '../components/ContainerList';

function App() {
  return (
    <div>

      <Navbar bg="dark" variant="dark">
        <Container>
          <Navbar.Brand>i5 Status</Navbar.Brand>
        </Container>
      </Navbar>

      <Container className="mt-3">
        <ContainerList />
      </Container>

    </div>
  );
}

export default App;
