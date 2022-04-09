import axios from 'axios';
import { useCallback, useEffect, useState } from 'react';

const Containers = () => {

  const [containers, setContainers] = useState([]);

  const reload = useCallback(
    async () => {
      try {
        setContainers(
          (await axios.get('/api/containers')).data
        );
      } catch (e) {

        // TODO: display error
      }
    },
    []
  );

  useEffect(() => { reload(); }, [reload]);

  return (
    <div>
      <h1>Containers</h1>
      <p>
        The table below lists all of the containers that are currently running in Docker.
      </p>

      <table className="table table-striped">
        <thead>
          <tr>
            <th>Name</th>
            <th>Domain</th>
            <th>{/* Status */}</th>
          </tr>
        </thead>
        <tbody>
          {containers.map(container => (
            <tr key={container.name}>
              <td>
                <a href={"http://" + container.domain}>
                  {container.name}
                </a>
              </td>
              <td className="text-muted">{container.domain}</td>
              <td className="text-end">
                {container.running ?
                  <span className="badge bg-success">Running</span> :
                  <span className="badge bg-danger">Stopped</span>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Containers;
