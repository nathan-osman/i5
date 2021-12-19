import React from 'react';
import axios from 'axios';
import { Table } from 'react-bootstrap';

export default class ContainerList extends React.Component {
  state = {
    containers: []
  }

  componentDidMount() {
    axios.get('/api/containers')
      .then(res => {
        this.setState({ containers: res.data })
      })
  }

  render() {
    return (
      <div>
        <h1>Containers</h1>
        <p>
          The table below lists all of the containers that are currently running in Docker.
        </p>

        <Table striped>
          <thead>
            <th>Name</th>
            <th>{/* Status */}</th>
          </thead>
          <tbody>
            {this.state.containers.map(container => (
              <tr>
                <td>
                  <a href={"http://" + container.domain}>
                    {container.name}
                  </a>
                </td>
                <td className="text-end">
                  {container.running ?
                    <span class="badge bg-success">Running</span> :
                    <span class="badge bg-danger">Stopped</span>}
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </div>
    )
  }
}
