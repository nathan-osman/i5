import { useEffect, useState } from "react";

const Requests = () => {

  const [requests, setRequests] = useState([]);

  useEffect(() => {
    const loc = window.location;
    let hostname = process.env.REACT_APP_HOSTNAME || loc.host;
    let websocket = new WebSocket(
      `${loc.protocol === "https:" ? 'wss:' : 'ws:'}${hostname}/api/ws`
    );
    websocket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      const date = new Date();
      const time =
        <>
          {date.getHours().toString().padStart(2, '0')}:
          {date.getMinutes().toString().padStart(2, '0')}:
          {date.getSeconds().toString().padStart(2, '0')}
          <small className="text-muted text-small">
            .{date.getMilliseconds().toString().padStart(3, '0')}
          </small>
        </>;
      setRequests(r => [
        {
          ...data.data,
          date,
          time
        },
        ...r
      ]);
    };
    return () => {
      websocket.close();
    };
  }, []);

  return (
    <div>
      <h1>Requests</h1>
      <p>
        The table below shows requests in realtime as they come in.
      </p>

      <table className="table table-striped">
        <thead>
          <tr>
            <th>Time</th>
            <th>Client</th>
            <th>Method</th>
            <th>Path</th>
          </tr>
        </thead>
        <tbody>
          {requests.length ?
            requests.map(request => (
              <tr key={+ request.date}>
                <td>{request.time}</td>
                <td>{request.remote_addr}</td>
                <td>
                  <span className="badge bg-secondary">{request.method}</span></td>
                <td><strong>{request.host}</strong>{request.path}</td>
              </tr>
            )) :
            <tr>
              <td colSpan="4">
                <p className="text-muted text-center p-5">No requests yet</p>
              </td>
            </tr>
          }
        </tbody>
      </table>
    </div>
  );
};

export default Requests;
