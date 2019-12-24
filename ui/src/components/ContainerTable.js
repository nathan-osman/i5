import React from 'react';
import './ContainerTable.scss';

const ContainerTable = ({ containers }) =>
  <table className="ContainerTable">
    <thead>
      <tr>
        <th>Name</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      {containers.map(container => (
        <tr>
          <td>{container.name}</td>
          <td>{container.status}</td>
        </tr>
      ))}
    </tbody>
  </table>;

export default ContainerTable;
