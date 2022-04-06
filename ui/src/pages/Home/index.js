import ContainerList from '../../components/ContainerList';
import RequestTicker from '../../components/RequestTicker';

const Home = () => {
  return (
    <div className="row">
      <div className="col-4">
        <ContainerList />
      </div>
      <div className="col-8">
        <RequestTicker />
      </div>
    </div>
  );
};

export default Home;
