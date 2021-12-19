import { Container, Navbar } from 'react-bootstrap';
import ContainerList from './components/ContainerList';

function App() {
  return (
    <div>

      <Navbar>
        <Container>
          <Navbar.Brand>i5 Status</Navbar.Brand>
        </Container>
      </Navbar>

      <Container>
        <ContainerList />
      </Container>

    </div>
  );
}

export default App;
