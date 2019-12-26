import React from 'react';
import axios from 'axios';
import './ContainerList.scss';

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
      <div className="ContainerList">
        <div className="header">
          <div className="cell name">Name</div>
          <div className="cell status">Status</div>
        </div>
        <div className="body">
          {this.state.containers.map(container => (
            <div className="row">
              <div className="cell name">{container.name}</div>
              {container.running ?
                <div className="cell status running">RUNNING</div> :
                <div className="cell status stopped">STOPPED</div>}
            </div>
          ))}
        </div>
      </div>
    )
  }
}
